package controllers

import (
	"net/http"

	"github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/viewparams/table"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/services/search"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewUploadSearchHandler(r router.Router) http.HandlerFunc {
	searchService := r.Context().Service(tokens.UploadSearchServiceToken).(search.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		searchRequest, err := request.GetSearchRequestDto(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get search request dto")
			r.Renderer().InternalServerError(res)
			return
		}

		s, err := request.MapSearchRequest(searchRequest)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request dto")
			r.Renderer().BadRequest(res)
			return
		}

		pt := new(entities.Upload)
		uploadResultLoader := func() ([]interface{}, int, error) {
			return searchService.Search(pt, searchRequest.Query, searchRequest.Count, searchRequest.Offset)
		}

		params, err := table.New(constants.UploadEntityName, pt, s, uploadResultLoader)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to build table params")
			r.Renderer().InternalServerError(res)
			return

		}

		r.Renderer().TableView(res, "uploadsdatatable", params)
	}
}
