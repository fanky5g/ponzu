package api

import (
	"github.com/fanky5g/ponzu/internal/application/content"
	"github.com/fanky5g/ponzu/internal/application/storage"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"log"
	"net/http"
)

func NewContentHandler(contentService content.Service, storageService storage.Service) http.HandlerFunc {
	handleCreateContent := NewCreateContentHandler(contentService, storageService)
	handleListContent := NewListContentHandler(contentService)
	handleGetContentById := NewContentByIdHandler(contentService)
	handleGetContentBySlug := NewContentBySlugHandler(contentService)
	handleUpdateContent := NewUpdateContentHandler(contentService, storageService)

	return func(res http.ResponseWriter, req *http.Request) {
		isSlug, identifier := request.GetRequestContentId(req)

		switch req.Method {
		case http.MethodPost:
			fallthrough
		case http.MethodPut:
			if identifier != "" {
				handleUpdateContent(res, req)
				return
			}

			handleCreateContent(res, req)
			return
		case http.MethodGet:
			if identifier != "" {
				if isSlug {
					handleGetContentBySlug(identifier, res, req)
					return
				}

				handleGetContentById(identifier, res, req)
				return
			}

			handleListContent(res, req)
			return
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}

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
		p, err := hook.BeforeAPIResponse(res, req, posts)
		if err != nil {
			log.Println("[Response] error calling BeforeAPIResponse:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		posts = p.([]interface{})
		writeJSONData(res, http.StatusOK, posts)

		// hook after response
		err = hook.AfterAPIResponse(res, req, posts)
		if err != nil {
			log.Println("[Response] error calling AfterAPIResponse:", err)
			return
		}
	}
}

func NewContentByIdHandler(contentService content.Service) func(string, http.ResponseWriter, *http.Request) {
	return func(contentId string, res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")

		pt, ok := item.Types[t]
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		post, err := contentService.GetContent(t, contentId)
		if err != nil {
			log.Printf("Failed to get content: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// assert hookable
		get := pt()
		hook, ok := get.(item.Hookable)
		if !ok {
			log.Println("[Response] error: Type", t, "does not implement item.Hookable or embed item.Item.")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// hook before response
		post, err = hook.BeforeAPIResponse(res, req, post)
		if err != nil {
			log.Println("[Response] error calling BeforeAPIResponse:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSONData(res, http.StatusOK, post)

		// hook after response
		err = hook.AfterAPIResponse(res, req, post)
		if err != nil {
			log.Println("[Response] error calling AfterAPIResponse:", err)
			return
		}
	}
}

func NewContentBySlugHandler(contentService content.Service) func(string, http.ResponseWriter, *http.Request) {
	return func(contentId string, res http.ResponseWriter, req *http.Request) {
		t, post, err := contentService.GetContentBySlug(contentId)
		if err != nil {
			log.Printf("Failed to get content: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if post == nil {
			writeJSONData(res, http.StatusNotFound, nil)
			return
		}

		pt, ok := item.Types[t]
		if !ok {
			writeJSONError(res, http.StatusBadRequest, item.ErrTypeNotRegistered)
			return
		}

		// assert hookable
		get := pt()
		hook, ok := get.(item.Hookable)
		if !ok {
			log.Println("[Response] error: Type", t, "does not implement item.Hookable or embed item.Item.")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// hook before response
		post, err = hook.BeforeAPIResponse(res, req, post)
		if err != nil {
			log.Println("[Response] error calling BeforeAPIResponse:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSONData(res, http.StatusOK, post)

		// hook after response
		err = hook.AfterAPIResponse(res, req, post)
		if err != nil {
			log.Println("[Response] error calling AfterAPIResponse:", err)
			return
		}
	}
}
