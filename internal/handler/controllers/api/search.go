package api

import (
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/services/search"
	"log"
	"net/http"
)

var ErrMissingSearchQuery = errors.New("query cannot be empty")

func NewSearchContentHandler(searchService search.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		t := req.URL.Query().Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		_, ok := item.Types[t]
		if !ok {
			writeJSONError(res, http.StatusBadRequest, fmt.Errorf(item.ErrTypeNotRegistered.Error(), t))
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
			log.Printf("[Search] error: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSONData(res, http.StatusOK, matches)
	}
}
