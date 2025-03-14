package content

import (
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/search"
	"log"
	"net/http"
)

var ErrMissingSearchQuery = errors.New("query cannot be empty")

func NewAPISearchContentHandler(contentTypes map[string]content.Builder, searchService *search.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		t := req.URL.Query().Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		entity, ok := contentTypes[t]
		if !ok {
			response.Respond(
				res,
				req,
				response.NewJSONResponse(
					http.StatusBadRequest,
					nil,
					fmt.Errorf(content.ErrTypeNotRegistered.Error(), t),
				),
			)
			return
		}

		searchRequest, err := search.GetSearchRequestDto(req)
		if err != nil {
			response.Respond(
				res,
				req,
				response.NewJSONResponse(
					http.StatusBadRequest,
					nil,
					err,
				),
			)
			return
		}

		if searchRequest.Query == "" {
			response.Respond(
				res,
				req,
				response.NewJSONResponse(
					http.StatusBadRequest,
					nil,
					ErrMissingSearchQuery,
				),
			)
			return
		}

		// TODO: implement pagination using response size
		matches, _, err := searchService.Search(entity(), searchRequest.Query, searchRequest.Count, searchRequest.Offset)
		if err != nil {
			log.Printf("[Find] error: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewJSONResponse(
				http.StatusOK,
				matches,
				nil,
			),
		)
	}
}
