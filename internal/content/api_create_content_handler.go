package content

import (
	"fmt"
	"github.com/fanky5g/ponzu/internal/http/response"
	"net/http"

	contentPkg "github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/http/request"
	log "github.com/sirupsen/logrus"
)

func NewCreateContentHandler(contentService *Service, contentTypes map[string]contentPkg.Builder) http.HandlerFunc {
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
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := contentService.CreateContent(t, post)
		if err != nil {
			log.Println("[Create] error:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// create JSON response to send data back to client
		response.Respond(
			res,
			req,
			response.NewJSONResponse(http.StatusOK, map[string]interface{}{
				"id":   id,
				"type": t,
			}, nil),
		)
	}
}
