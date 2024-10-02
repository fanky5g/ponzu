package controllers

import (
	"net/http"
	"strings"

	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewDeleteHandler(r router.Router) http.HandlerFunc {
	contentService := r.Context().Service(tokens.ContentServiceToken).(content.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse form")
			r.Renderer().InternalServerError(res)
			return
		}

		ids := make([]string, 0)

		idParam := strings.TrimSpace(req.FormValue("id"))
		idsParam := strings.TrimSpace(req.FormValue("ids"))
		if idParam != "" {
			ids = append(ids, idParam)
		} else if idsParam != "" {
			idsToDelete := strings.FieldsFunc(idsParam, func(c rune) bool {
				return c == ','
			})

			for _, idToDelete := range idsToDelete {
				ids = append(ids, strings.TrimSpace(idToDelete))
			}
		}

		t := req.FormValue("type")

		if len(ids) == 0 || t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = contentService.DeleteContent(t, ids...)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to delete content")
			r.Renderer().InternalServerError(res)
			return
		}

		r.Redirect(req, res, "/contents?type="+t)
	}
}
