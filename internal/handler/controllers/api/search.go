package api

import (
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/search"
	"github.com/fanky5g/ponzu/tokens"
	"log"
	"net/http"
)

var ErrMissingSearchQuery = errors.New("query cannot be empty")

func NewSearchContentHandler(r router.Router) http.HandlerFunc {
	contentTypes := r.Context().Types().Content
	searchService := r.Context().Service(tokens.ContentSearchServiceToken).(search.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		t := req.URL.Query().Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		_, ok := contentTypes[t]
		if !ok {
			r.Renderer().Error(res, http.StatusBadRequest, fmt.Errorf(content.ErrTypeNotRegistered.Error(), t))
			return
		}

		searchRequest, err := request.GetSearchRequestDto(req)
		if err != nil {
			r.Renderer().Error(res, http.StatusBadRequest, err)
			return
		}

		if searchRequest.Query == "" {
			r.Renderer().Error(res, http.StatusBadRequest, ErrMissingSearchQuery)
			return
		}

		// TODO: implement pagination using response size
		matches, _, err := searchService.Search(t, searchRequest.Query, searchRequest.Count, searchRequest.Offset)
		if err != nil {
			log.Printf("[Find] error: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Renderer().Json(res, http.StatusOK, matches)
	}
}
