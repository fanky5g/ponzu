package controllers

import (
	"github.com/fanky5g/ponzu/internal/content"
	"net/http"

	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/viewparams/table"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewUploadContentsHandler(r router.Router) http.HandlerFunc {
	uploadService := r.Context().Service(tokens.UploadServiceToken).(*content.UploadService)

	return func(res http.ResponseWriter, req *http.Request) {
		pt := new(entities.Upload)
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
			return uploadService.GetAllWithOptions(search)
		}

		params, err := table.New(constants.UploadEntityName, pt, search, uploadResultLoader)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to build table params")
			r.Renderer().InternalServerError(res)
			return

		}

		r.Renderer().TableView(res, "templates/uploadsdatatable", params)
	}
}
