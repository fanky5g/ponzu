package api

import (
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/internal/application/content"
	"github.com/fanky5g/ponzu/internal/application/storage"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"log"
	"net/http"
)

func NewUpdateContentHandler(contentService content.Service, storageService storage.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		isSlug, identifier := request.GetRequestContentId(req)
		if identifier == "" {
			writeJSONError(res, http.StatusBadRequest, errors.New("content id is required"))
			return
		}

		if isSlug {
			writeJSONError(res, http.StatusBadRequest, errors.New("slug not supported for update"))
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

		pt, ok := item.Types[t]
		if !ok {
			writeJSONError(res, http.StatusBadRequest, fmt.Errorf(item.ErrTypeNotRegistered.Error(), t))
			return
		}

		hook, ok := pt().(item.Hookable)
		if !ok {
			log.Println("[Update] error: Type", t, "does not implement item.Hookable or embed item.Item.")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		update, err := request.MapRequestToContentUpdate(req)
		if err != nil {
			writeJSONError(res, http.StatusBadRequest, err)
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

		writeJSONData(res, http.StatusOK, u)
	}
}
