package api

import (
	"github.com/fanky5g/ponzu/internal/application/content"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"log"
	"net/http"
)

func NewListContentHandler(contentService content.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		it, ok := item.Types[t]
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		searchRequestDto, err := request.GetSearchRequestDto(req.URL.Query())
		if err != nil {
			writeJSONError(res, http.StatusBadRequest, err)
			return
		}

		search, err := request.MapSearchRequest(searchRequestDto)
		if err != nil {
			writeJSONError(res, http.StatusBadRequest, err)
			return
		}

		_, posts, err := contentService.GetAllWithOptions(t, search)
		if err != nil {
			log.Printf("Failed to list content: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// assert hookable
		get := it()
		hook, ok := get.(item.Hookable)
		if !ok {
			log.Println("[Response] error: Type", t, "does not implement item.Hookable or embed item.Item.")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// hook before response
		posts, err = hook.BeforeAPIResponse(res, req, posts)
		if err != nil {
			log.Println("[Response] error calling BeforeAPIResponse:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSONData(res, http.StatusOK, posts)

		// hook after response
		err = hook.AfterAPIResponse(res, req, posts)
		if err != nil {
			log.Println("[Response] error calling AfterAPIResponse:", err)
			return
		}
	}
}
