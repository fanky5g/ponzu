package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	contentPkg "github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/tokens"
)

func NewUpdateContentHandler(r router.Router) http.HandlerFunc {
	contentTypes := r.Context().Types().Content
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)
	uploadService := r.Context().Service(tokens.UploadServiceToken).(*content.UploadService)

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
			savedFiles, err := uploadService.UploadFiles(files)
			if err != nil {
				log.Println("[Update] error:", err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			for _, file := range savedFiles {
				req.PostForm.Set(file.Name, file.Path)
			}
		}

		pt, ok := contentTypes[t]
		if !ok {
			r.Renderer().Error(res, http.StatusBadRequest, fmt.Errorf(contentPkg.ErrTypeNotRegistered.Error(), t))
			return
		}

		update, err := request.GetEntity(pt, req)
		if err != nil {
			r.Renderer().Error(res, http.StatusBadRequest, err)
			return
		}

		u, err := contentService.UpdateContent(t, identifier, update)
		if err != nil {
			log.Printf("[Update] error: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Renderer().Json(res, http.StatusOK, u)
	}
}
