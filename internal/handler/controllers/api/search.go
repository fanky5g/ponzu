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
			writeJSONError(res, http.StatusBadRequest, fmt.Errorf(content.ErrTypeNotRegistered.Error(), t))
			return
		}

		searchRequest, err := request.GetSearchRequestDto(req)
		if err != nil {
			writeJSONError(res, http.StatusBadRequest, err)
			return
		}

		if searchRequest.Query == "" {
			writeJSONError(res, http.StatusBadRequest, ErrMissingSearchQuery)
			return
		}

		matches, err := searchService.Search(t, searchRequest.Query, searchRequest.Count, searchRequest.Offset)
		if err != nil {
			log.Printf("[Find] error: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSONData(res, http.StatusOK, matches)
	}
}
