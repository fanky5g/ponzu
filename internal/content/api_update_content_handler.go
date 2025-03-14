package content

import (
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/internal/http/response"
	"log"
	"net/http"

	contentPkg "github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/http/request"
)

func NewUpdateContentHandler(
	contentTypes map[string]contentPkg.Builder,
	contentService *Service,
	uploadService *UploadService,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		isSlug, identifier := request.GetRequestContentId(req)
		if identifier == "" {
			response.Respond(
				res,
				req,
				response.NewJSONResponse(
					http.StatusBadRequest,
					nil,
					errors.New("entities id is required"),
				),
			)
			return
		}

		if isSlug {
			response.Respond(
				res,
				req,
				response.NewJSONResponse(
					http.StatusBadRequest,
					nil,
					errors.New("slug not supported for update"),
				),
			)
			return
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
			response.Respond(
				res,
				req,
				response.NewJSONResponse(
					http.StatusBadRequest,
					nil,
					fmt.Errorf(contentPkg.ErrTypeNotRegistered.Error(), t),
				),
			)
			return
		}

		update, err := request.GetEntity(pt, req)
		if err != nil {
			response.Respond(
				res,
				req,
				response.NewJSONResponse(
					http.StatusBadRequest,
					nil,
					err,
				),
			)
			return
		}

		u, err := contentService.UpdateContent(t, identifier, update)
		if err != nil {
			log.Printf("[Update] error: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewJSONResponse(
				http.StatusOK,
				u,
				nil,
			),
		)
	}
}
