package content

import (
	"net/http"
	"strings"

	"github.com/fanky5g/ponzu/internal/http/response"
	log "github.com/sirupsen/logrus"
)

func NewExportHandler(contentService *Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")
		format := strings.ToLower(q.Get("format"))

		if t == "" || format == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		ds, err := contentService.Export(t, format)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to export")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(res, req, response.NewStreamResponse(http.StatusOK, ds))
	}
}
