package controllers

import (
	"net/http"

	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewDeleteHandler(r router.Router) http.HandlerFunc {
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		selectedItems, err := request.GetSelectedItems(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse request")
			r.Renderer().InternalServerError(res)
		}

		if len(selectedItems) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		t := req.FormValue("type")

		if len(selectedItems) == 0 || t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = contentService.DeleteContent(t, selectedItems...)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to delete content")
			r.Renderer().InternalServerError(res)
			return
		}

		r.Redirect(req, res, "/contents?type="+t)
	}
}
