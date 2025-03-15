package content

import (
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	"github.com/fanky5g/ponzu/internal/search"
	"net/http"

	"github.com/fanky5g/ponzu/content/entities"
	log "github.com/sirupsen/logrus"
)

func NewUploadSearchHandler(publicPath string, searchService *search.Service, tmpl layouts.Template) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		searchRequest, err := search.GetSearchRequestDto(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get search request dto")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		s, err := search.MapSearchRequest(searchRequest)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map search request dto")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		pt := new(entities.Upload)
		uploadResultLoader := func() ([]interface{}, int, error) {
			return searchService.Search(pt, searchRequest.Query, searchRequest.Count, searchRequest.Offset)
		}

		data, err := buildTableViewModel(publicPath, pt, s, uploadResultLoader)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to build table params")
			res.WriteHeader(http.StatusInternalServerError)
			return

		}

		response.Respond(
			res,
			req,
			response.NewHTMLResponse(
				http.StatusOK,
				tmpl,
				data,
			),
		)
	}
}
