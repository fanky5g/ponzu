package content

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	"github.com/fanky5g/ponzu/internal/search"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewContentSearchHandler(
	publicPath string,
	searchService *search.Service,
	contentTypes map[string]content.Builder,
	dataTable layouts.Template,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")

		searchRequest, err := search.GetSearchRequestDto(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request DTO")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Query must be set
		if searchRequest.Query == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if t == "" {
			response.Respond(
				res,
				req,
				response.NewRedirectResponse(publicPath, "/"),
			)
			return
		}

		pt, ok := contentTypes[t]
		if !ok {
			response.Respond(
				res,
				req,
				response.NewRedirectResponse(publicPath, "/"),
			)
			return
		}

		entity, ok := pt().(content.Entity)
		if !ok {
			log.Warnf("item %s does not implement entity interface", t)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		s, err := search.MapSearchRequest(searchRequest)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		searchResultsLoader := func() ([]interface{}, int, error) {
			return searchService.Search(
				pt(),
				searchRequest.Query,
				searchRequest.Count,
				searchRequest.Offset,
			)
		}

		data, err := buildTableViewModel(publicPath, entity, s, searchResultsLoader)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to build table view model")
			res.WriteHeader(http.StatusInternalServerError)
			return

		}

		response.Respond(
			res,
			req,
			response.NewHTMLResponse(
				http.StatusOK,
				dataTable,
				data,
			),
		)
	}
}
