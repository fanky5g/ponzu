package controllers

import (
	"net/http"

	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/viewparams/table"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewUploadContentsHandler(r router.Router) http.HandlerFunc {
	storageService := r.Context().Service(tokens.StorageServiceToken).(storage.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		pt := new(entities.FileUpload)
		_, ok := interface{}(pt).(editor.Editable)
		if !ok {
			log.Warning("entities.FileUpload is not editable")
			r.Renderer().InternalServerError(res)
			return
		}

		searchRequestDto, err := request.GetSearchRequestDto(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get search request dto")
			return
		}

		search, err := request.MapSearchRequest(searchRequestDto)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request dto")
			return
		}

		uploadResultLoader := func() ([]interface{}, int, error) {
			total, matches, err := storageService.GetAllWithOptions(search)
			if err != nil {
				return nil, 0, err
			}

			if len(matches) > 0 {
				out := make([]interface{}, len(matches))
				for i := 0; i < len(matches); i++ {
					out[i] = matches[i]
				}

				return out, total, err
			}

			return nil, 0, nil
		}

		params, err := table.New(constants.UploadsEntityName, pt, search, uploadResultLoader)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to build table params")
			r.Renderer().InternalServerError(res)
			return

		}

		r.Renderer().TableView(res, "templates/uploadsdatatable", params)
	}
}
