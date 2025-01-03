package api

import (
	"context"
	"fmt"
	"net/http"

	contentPkg "github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewCreateContentHandler(r router.Router) http.HandlerFunc {
	contentTypes := r.Context().Types().Content
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)
	storageService := r.Context().Service(tokens.StorageServiceToken).(storage.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		t := req.URL.Query().Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		contentType, ok := contentTypes[t]
		if !ok {
			_, err := fmt.Fprintf(res, contentPkg.ErrTypeNotRegistered.Error(), t)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to write response")
			}

			return
		}

		files, err := request.GetRequestFiles(req)
		if err != nil {
			log.Println("[Create] error:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(files) > 0 {
			var urlPaths map[string]string
			urlPaths, err = storageService.UploadFiles(files)
			if err != nil {
				log.Println("[Create] error:", err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			for name, urlPath := range urlPaths {
				req.PostForm.Set(name, urlPath)
			}
		}

		post, err := request.GetEntity(contentType, req)
		if err != nil {
			r.Renderer().Error(res, http.StatusBadRequest, err)
			return
		}

		hook, ok := post.(item.Hookable)
		if !ok {
			log.Println("[Create] error: Type", t, "does not implement item.Hookable or embed item.Item.")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = hook.BeforeAPICreate(res, req)
		if err != nil {
			log.Println("[Create] error calling BeforeCreate:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = hook.BeforeSave(res, req)
		if err != nil {
			log.Println("[Create] error calling BeforeSave:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		id, err := contentService.CreateContent(t, post)
		if err != nil {
			log.Println("[Create] error:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// set the target in the context so user can get saved value in hooks
		ctx := context.WithValue(req.Context(), "target", fmt.Sprintf("%s:%s", t, id))
		req = req.WithContext(ctx)

		err = hook.AfterSave(res, req)
		if err != nil {
			log.Println("[Create] error calling AfterSave:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = hook.AfterAPICreate(res, req)
		if err != nil {
			log.Println("[Create] error calling AfterAccept:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// create JSON response to send data back to client
		r.Renderer().Json(res, http.StatusOK, map[string]interface{}{
			"id":   id,
			"type": t,
		})
	}
}
