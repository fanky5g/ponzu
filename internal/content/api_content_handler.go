package content

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/search"
	"log"
	"net/http"
)

func NewAPIContentHandler(contentService *Service, uploadService *UploadService, contentTypes map[string]content.Builder) http.HandlerFunc {
	handleCreateContent := NewCreateContentHandler(contentService, contentTypes)
	handleListContent := NewListContentHandler(contentService, contentTypes)
	handleGetContentById := NewContentByIdHandler(contentService, contentTypes)
	handleGetContentBySlug := NewContentBySlugHandler(contentService)
	handleUpdateContent := NewUpdateContentHandler(contentTypes, contentService, uploadService)
	handleDeleteContent := NewAPIDeleteContentHandler(contentTypes, contentService)

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

func NewListContentHandler(contentService *Service, contentTypes map[string]content.Builder) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		_, ok := contentTypes[t]
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		searchRequestDto, err := search.GetSearchRequestDto(req)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		s, err := search.MapSearchRequest(searchRequestDto)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		posts, _, err := contentService.GetAllWithOptions(t, s)
		if err != nil {
			log.Printf("Failed to list entities: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewJSONResponse(http.StatusOK, posts, nil),
		)
	}
}

func NewContentByIdHandler(
	contentService *Service,
	contentTypes map[string]content.Builder,
) func(string, http.ResponseWriter, *http.Request) {
	return func(contentId string, res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		t := q.Get("type")

		_, ok := contentTypes[t]
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

		response.Respond(
			res,
			req,
			response.NewJSONResponse(http.StatusOK, post, nil),
		)
	}
}

func NewContentBySlugHandler(contentService *Service) func(string, http.ResponseWriter, *http.Request) {
	return func(contentId string, res http.ResponseWriter, req *http.Request) {
		post, err := contentService.GetContentBySlug(contentId)
		if err != nil {
			log.Printf("Failed to get entities: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if post == nil {
			response.Respond(
				res,
				req,
				response.NewJSONResponse(http.StatusNotFound, nil, nil),
			)
			return
		}

		response.Respond(
			res,
			req,
			response.NewJSONResponse(http.StatusOK, post, nil),
		)
	}
}
