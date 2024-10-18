package controllers

import (
	"net/http"

	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/viewparams/table"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/search"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewSearchHandler(r router.Router) http.HandlerFunc {
	searchService := r.Context().Service(tokens.ContentSearchServiceToken).(search.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")

		searchRequest, err := request.GetSearchRequestDto(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request DTO")
			r.Renderer().InternalServerError(res)
			return
		}

		// Query must be set
		if searchRequest.Query == "" {
			r.Renderer().BadRequest(res)
			return
		}

		if t == "" {
			r.Redirect(req, res, "/admin")
			return
		}

		pt, ok := r.Context().Types().Content[t]
		if !ok {
			r.Redirect(req, res, "/admin")
			return
		}

		search, err := request.MapSearchRequest(searchRequest)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request")
			r.Renderer().InternalServerError(res)
			return
		}

		searchResultsLoader := func() ([]interface{}, int, error) {
			return searchService.Search(
				pt(),
				// TODO: refactor search service to accept only one searchRequest argument
				searchRequest.Query,
				searchRequest.Count,
				searchRequest.Offset,
			)
		}

		params, err := table.New(t, pt, search, searchResultsLoader)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to build table params")
			r.Renderer().InternalServerError(res)
			return

		}

		r.Renderer().TableView(res, "templates/datatable", params)
	}
}
