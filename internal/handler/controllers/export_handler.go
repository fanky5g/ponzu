package controllers

import (
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

func NewExportHandler(r router.Router) http.HandlerFunc {
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		// /contents/export?type=Blogpost&format=csv
		q := req.URL.Query()
		t := q.Get("type")
		f := strings.ToLower(q.Get("format"))

		if t == "" || f == "" {
			r.Renderer().BadRequest(res)
			return
		}

		switch f {
		case "csv":
			response, err := contentService.ExportCSV(t)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to export")
				r.Renderer().InternalServerError(res)
				return
			}

			if response == nil {
				res.WriteHeader(http.StatusNoContent)
				return
			}

			res.Header().Set("Content-Type", response.ContentType)
			res.Header().Set("Content-Disposition", response.ContentDisposition)
			if _, err = io.Copy(res, response.Payload); err != nil {
				log.WithField("Error", err).Warning("Failed to write response")
			}
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
