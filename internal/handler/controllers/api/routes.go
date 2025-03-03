package api

import (
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/storage"
	"net/http"
)

func RegisterRoutes(r router.Router, uploadsFileSystem storage.Client) {
	r.APIRoute("/api/auth", NewAuthHandler)
	r.APIAuthorizedRoute("/api/content/", NewContentHandler)
	r.APIAuthorizedRoute("/api/search", NewSearchContentHandler)

	r.HandleWithCacheControl(
		"/api/uploads/",
		http.StripPrefix("/api/uploads/", http.FileServer(uploadsFileSystem)),
	)
}
