package api

import (
	"context"
	"fmt"
	"github.com/fanky5g/ponzu/internal/application/content"
	"github.com/fanky5g/ponzu/internal/application/storage"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"log"
	"net/http"
)

func NewCreateContentHandler(contentService content.Service, storageService storage.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		t := req.URL.Query().Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
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
			urlPaths, err = storageService.StoreFiles(files)
			if err != nil {
				log.Println("[Create] error:", err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			for name, urlPath := range urlPaths {
				req.PostForm.Set(name, urlPath)
			}
		}

		post, err := request.GetEntity(t, req)
		if err != nil {
			writeJSONError(res, http.StatusBadRequest, err)
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
			return
		}

		err = hook.BeforeSave(res, req)
		if err != nil {
			log.Println("[Create] error calling BeforeSave:", err)
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
			return
		}

		err = hook.AfterAPICreate(res, req)
		if err != nil {
			log.Println("[Create] error calling AfterAccept:", err)
			return
		}

		// create JSON response to send data back to client
		writeJSONData(res, http.StatusOK, map[string]interface{}{
			"id":   id,
			"type": t,
		})
	}
}
