package api

import (
	"errors"
	"fmt"
	contentPkg "github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/tokens"
	"log"
	"net/http"
)

func NewUpdateContentHandler(r router.Router) http.HandlerFunc {
	contentTypes := r.Context().Types().Content
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)
	storageService := r.Context().Service(tokens.StorageServiceToken).(storage.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		isSlug, identifier := request.GetRequestContentId(req)
		if identifier == "" {
			r.Renderer().Error(res, http.StatusBadRequest, errors.New("entities id is required"))
			return
		}

		if isSlug {
			r.Renderer().Error(res, http.StatusBadRequest, errors.New("slug not supported for update"))
		}

		t := req.URL.Query().Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		files, err := request.GetRequestFiles(req)
		if err != nil {
			log.Println("[Update] error:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(files) > 0 {
			var urlPaths map[string]string
			urlPaths, err = storageService.StoreFiles(files)
			if err != nil {
				log.Println("[Update] error:", err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			for name, urlPath := range urlPaths {
				req.PostForm.Set(name, urlPath)
			}
		}

		pt, ok := contentTypes[t]
		if !ok {
			r.Renderer().Error(res, http.StatusBadRequest, fmt.Errorf(contentPkg.ErrTypeNotRegistered.Error(), t))
			return
		}

		hook, ok := pt().(item.Hookable)
		if !ok {
			log.Println("[Update] error: Type", t, "does not implement item.Hookable or embed item.Item.")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		update, err := request.GetEntity(pt, req)
		if err != nil {
			r.Renderer().Error(res, http.StatusBadRequest, err)
			return
		}

		err = hook.BeforeAPIUpdate(res, req)
		if err != nil {
			log.Println("[Update] error calling BeforeUpdate:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = hook.BeforeSave(res, req)
		if err != nil {
			log.Println("[Create] error calling BeforeSave:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		u, err := contentService.UpdateContent(t, identifier, update)
		if err != nil {
			log.Printf("[Update] error: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = hook.AfterSave(res, req)
		if err != nil {
			log.Println("[Create] error calling AfterSave:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = hook.AfterAPIUpdate(res, req)
		if err != nil {
			log.Println("[Update] error calling AfterUpdate:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Renderer().Json(res, http.StatusOK, u)
	}
}
