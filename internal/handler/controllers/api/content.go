package api

import (
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/tokens"
	"log"
	"net/http"
)

func NewContentHandler(r router.Router) http.HandlerFunc {
	handleCreateContent := NewCreateContentHandler(r)
	handleListContent := NewListContentHandler(r)
	handleGetContentById := NewContentByIdHandler(r)
	handleGetContentBySlug := NewContentBySlugHandler(r)
	handleUpdateContent := NewUpdateContentHandler(r)
	handleDeleteContent := NewDeleteContentHandler(r)

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
		case http.MethodDelete:
			handleDeleteContent(res, req)
			return
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}

func NewListContentHandler(r router.Router) http.HandlerFunc {
	contentTypes := r.Context().Types().Content
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		it, ok := contentTypes[t]
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		searchRequestDto, err := request.GetSearchRequestDto(req)
		if err != nil {
			r.Renderer().Error(res, http.StatusBadRequest, err)
			return
		}

		search, err := request.MapSearchRequest(searchRequestDto)
		if err != nil {
			r.Renderer().Error(res, http.StatusBadRequest, err)
			return
		}

		posts, _, err := contentService.GetAllWithOptions(t, search)
		if err != nil {
			log.Printf("Failed to list entities: %v\n", err)
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
		r.Renderer().Json(res, http.StatusOK, posts)

		// hook after response
		err = hook.AfterAPIResponse(res, req, posts)
		if err != nil {
			log.Println("[Response] error calling AfterAPIResponse:", err)
			return
		}
	}
}

func NewContentByIdHandler(r router.Router) func(string, http.ResponseWriter, *http.Request) {
	contentTypes := r.Context().Types().Content
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)

	return func(contentId string, res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")

		pt, ok := contentTypes[t]
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		post, err := contentService.GetContent(t, contentId)
		if err != nil {
			log.Printf("Failed to get entities: %v\n", err)
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

		r.Renderer().Json(res, http.StatusOK, post)

		// hook after response
		err = hook.AfterAPIResponse(res, req, post)
		if err != nil {
			log.Println("[Response] error calling AfterAPIResponse:", err)
			return
		}
	}
}

func NewContentBySlugHandler(r router.Router) func(string, http.ResponseWriter, *http.Request) {
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)

	return func(contentId string, res http.ResponseWriter, req *http.Request) {
		post, err := contentService.GetContentBySlug(contentId)
		if err != nil {
			log.Printf("Failed to get entities: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if post == nil {
			r.Renderer().Json(res, http.StatusNotFound, nil)
			return
		}

		// assert hookable
		hook, ok := post.(item.Hookable)
		if !ok {
			log.Println("[Response] error: Type does not implement item.Hookable or embed item.Item.")
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

		r.Renderer().Json(res, http.StatusOK, post)

		// hook after response
		err = hook.AfterAPIResponse(res, req, post)
		if err != nil {
			log.Println("[Response] error calling AfterAPIResponse:", err)
			return
		}
	}
}
