package controllers

import (
	"net/http"

	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewDeleteUploadHandler(r router.Router) http.HandlerFunc {
	storageService := r.Context().Service(tokens.StorageServiceToken).(storage.Service)

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

		post := interface{}(&entities.FileUpload{})
		hook, ok := post.(item.Hookable)
		if !ok {
			log.Println("Type", constants.UploadsEntityName, "does not implement item.Hookable or embed item.Item.")
			r.Renderer().BadRequest(res)
			return
		}

		err = hook.BeforeDelete(res, req)
		if err != nil {
			log.Println("Error running BeforeDelete method in deleteHandler for:", constants.UploadsEntityName, err)
			return
		}

		err = storageService.DeleteFile(selectedItems...)
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = hook.AfterDelete(res, req)
		if err != nil {
			log.Println("Error running AfterDelete method in deleteHandler for:", constants.UploadsEntityName, err)
			return
		}

		r.Redirect(req, res, "/uploads")
	}
}
