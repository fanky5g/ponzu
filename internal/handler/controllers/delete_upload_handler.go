package controllers

import (
	"net/http"

	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/uploads"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewDeleteUploadHandler(r router.Router) http.HandlerFunc {
	uploadService := r.Context().Service(tokens.UploadServiceToken).(*uploads.Service)

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

		err = uploadService.DeleteUpload(selectedItems...)
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Redirect(req, res, "/uploads")
	}
}
