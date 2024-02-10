package api

import (
	"github.com/fanky5g/ponzu/internal/application/content"
	"github.com/fanky5g/ponzu/internal/application/storage"
	"net/http"
)

func NewContentHandler(contentService content.Service, storageService storage.Service) http.HandlerFunc {
	handleCreateContent := NewCreateContentHandler(contentService, storageService)
	handleListContent := NewListContentHandler(contentService)

	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			handleCreateContent(res, req)
			return
		case http.MethodGet:
			handleListContent(res, req)
			return
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
