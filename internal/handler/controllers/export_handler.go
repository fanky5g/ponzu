package controllers

import (
	"net/http"
	"strings"

	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewExportHandler(r router.Router) http.HandlerFunc {
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")
		format := strings.ToLower(q.Get("format"))

		if t == "" || format == "" {
			r.Renderer().BadRequest(res)
			return
		}

		ds, err := contentService.Export(t, format)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to export")
			r.Renderer().InternalServerError(res)
			return
		}

		response.Write(res, response.NewStreamResponse(http.StatusOK, ds))
	}
}
