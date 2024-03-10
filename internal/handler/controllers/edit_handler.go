package controllers

import (
	"context"
	"fmt"
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/domain/services/content"
	"github.com/fanky5g/ponzu/internal/domain/services/management/editor"
	"github.com/fanky5g/ponzu/internal/domain/services/management/manager"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/views"
	"github.com/fanky5g/ponzu/internal/services/config"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/internal/util"
	"log"
	"net/http"
	"strings"
)

func NewEditHandler(
	pathConf conf.Paths,
	configService config.Service,
	contentService content.Service,
	storageService storage.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		appName, err := configService.GetAppName()
		if err != nil {
			log.Printf("Failed to get app name: %v\n", appName)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch req.Method {
		case http.MethodGet:
			q := req.URL.Query()
			i := q.Get("id")
			t := q.Get("type")
			status := q.Get("status")

			contentType, ok := item.Types[t]
			if !ok {
				fmt.Fprintf(res, item.ErrTypeNotRegistered.Error(), t)
				return
			}

			post := contentType()
			if i != "" {
				if status == "pending" {
					t = t + "__pending"
				}

				post, err = contentService.GetContent(t, i)
				if err != nil {
					LogAndFail(res, err, appName, pathConf)
					return
				}

				if post == nil {
					res.WriteHeader(http.StatusNotFound)
					errView, err := views.Admin(util.Html("error_404"), appName, pathConf)
					if err != nil {
						return
					}

					res.Write(errView)
					return
				}
			} else {
				_, ok = post.(item.Identifiable)
				if !ok {
					log.Println("Content type", t, "doesn't implement item.Identifiable")
					return
				}
			}

			m, err := manager.Manage(post.(editor.Editable), pathConf, t)
			if err != nil {
				LogAndFail(res, err, appName, pathConf)
				return
			}

			adminView, err := views.Admin(string(m), appName, pathConf)
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			res.Header().Set("Content-Type", "text/html")
			res.Write(adminView)
		case http.MethodPost:
			err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
			if err != nil {
				LogAndFail(res, err, appName, pathConf)
				return
			}

			cid := req.FormValue("id")
			t := req.FormValue("type")
			files, err := request.GetRequestFiles(req)
			if err != nil {
				LogAndFail(res, err, appName, pathConf)
				return
			}

			urlPaths, err := storageService.StoreFiles(files)
			if err != nil {
				LogAndFail(res, err, appName, pathConf)
				return
			}

			for name, urlPath := range urlPaths {
				req.PostForm.Set(name, urlPath)
			}

			pt := t
			if strings.Contains(t, "__") {
				pt = strings.Split(t, "__")[0]
			}

			entity, err := request.GetEntityFromFormData(pt, req.PostForm)
			if err != nil {
				LogAndFail(res, err, appName, pathConf)
				return
			}

			hook, ok := entity.(item.Hookable)
			if !ok {
				log.Println("Type", pt, "does not implement item.Hookable or embed item.Item.")
				res.WriteHeader(http.StatusBadRequest)
				errView, err := views.Admin(util.Html("error_400"), appName, pathConf)
				if err != nil {
					return
				}

				res.Write(errView)
				return
			}

			if cid == "" {
				err = hook.BeforeAdminCreate(res, req)
				if err != nil {
					log.Println("Error running BeforeAdminCreate method in editHandler for:", t, err)
					return
				}
			} else {
				err = hook.BeforeAdminUpdate(res, req)
				if err != nil {
					log.Println("Error running BeforeAdminUpdate method in editHandler for:", t, err)
					return
				}
			}

			err = hook.BeforeSave(res, req)
			if err != nil {
				log.Println("Error running BeforeSave method in editHandler for:", t, err)
				return
			}

			id, err := contentService.CreateContent(t, entity)
			if err != nil {
				LogAndFail(res, err, appName, pathConf)
				return
			}

			// set the target in the context so user can get saved value from db in hook
			ctx := context.WithValue(req.Context(), "target", fmt.Sprintf("%s:%s", t, id))
			req = req.WithContext(ctx)

			err = hook.AfterSave(res, req)
			if err != nil {
				log.Println("Error running AfterSave method in editHandler for:", t, err)
				return
			}

			if cid == "" {
				err = hook.AfterAdminCreate(res, req)
				if err != nil {
					log.Println("Error running AfterAdminUpdate method in editHandler for:", t, err)
					return
				}
			} else {
				err = hook.AfterAdminUpdate(res, req)
				if err != nil {
					log.Println("Error running AfterAdminUpdate method in editHandler for:", t, err)
					return
				}
			}

			path := strings.TrimSuffix(pathConf.PublicPath, req.URL.Path)
			redir := path + "?type=" + pt + "&id=" + id

			if req.URL.Query().Get("status") == "pending" {
				redir += "&status=pending"
			}

			util.Redirect(req, res, pathConf, redir, http.StatusFound)
			return
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
