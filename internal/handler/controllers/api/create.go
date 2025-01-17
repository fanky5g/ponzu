package api

import (
	"fmt"
	"net/http"

	contentPkg "github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewCreateContentHandler(r router.Router) http.HandlerFunc {
	contentTypes := r.Context().Types().Content
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)

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

		post, err := request.GetEntity(contentType, req)
		if err != nil {
			r.Renderer().Error(res, http.StatusBadRequest, err)
			return
		}

		id, err := contentService.CreateContent(t, post)
		if err != nil {
			log.Println("[Create] error:", err)
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
