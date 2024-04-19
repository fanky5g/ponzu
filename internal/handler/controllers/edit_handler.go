package controllers

import (
	"context"
	"fmt"
	contentPkg "github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewEditHandler(r router.Router) http.HandlerFunc {
	contentTypes := r.Context().Types().Content
	storageService := r.Context().Service(tokens.StorageServiceToken).(storage.Service)
	contentService := r.Context().Service(tokens.ContentServiceToken).(content.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			q := req.URL.Query()
			i := q.Get("id")
			t := q.Get("type")

			contentType, ok := contentTypes[t]
			if !ok {
				_, err := fmt.Fprintf(res, contentPkg.ErrTypeNotRegistered.Error(), t)
				if err != nil {
					log.WithField("Error", err).Warning("Failed to write response")
				}

				return
			}

			contentEntry := contentType()
			var err error
			if i != "" {
				contentEntry, err = contentService.GetContent(t, i)
				if err != nil {
					log.WithField("Error", err).Warning("Failed to get content")
					return
				}

				if contentEntry == nil {
					r.Renderer().BadRequest(res)
					return
				}
			} else {
				_, ok = contentEntry.(item.Identifiable)
				if !ok {
					log.Println("Content type", t, "doesn't implement item.Identifiable")
					return
				}
			}

			r.Renderer().ManageEditable(res, contentEntry.(editor.Editable), t)
		case http.MethodPost:
			cid := req.FormValue("id")
			t := req.FormValue("type")

			contentType, ok := contentTypes[t]
			if !ok {
				_, err := fmt.Fprintf(res, contentPkg.ErrTypeNotRegistered.Error(), t)
				if err != nil {
					log.WithField("Error", err).Warning("Failed to write response")
				}

				return
			}

			err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
			if err != nil {
				log.WithField("Error", err).Warning("Failed to parse form")
				r.Renderer().InternalServerError(res)
				return
			}

			files, err := request.GetRequestFiles(req)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to get request files")
				r.Renderer().InternalServerError(res)
				return
			}

			urlPaths, err := storageService.StoreFiles(files)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to get save files")
				return
			}

			for name, urlPath := range urlPaths {
				req.PostForm.Set(name, urlPath)
			}

			entity, err := request.GetEntityFromFormData(contentType, req.PostForm)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to map request entity")
				return
			}

			hook, ok := entity.(item.Hookable)
			if !ok {
				log.Println("Type", t, "does not implement item.Hookable or embed item.Item.")
				r.Renderer().BadRequest(res)
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
				log.WithField("Error", err).Warning("Failed to create content")
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

			r.Redirect(req, res, req.URL.Path+"?type="+t+"&id="+id)
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
