package api

import (
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"net/http"
)

func RegisterRoutes(r router.Router, uploadsFileSystem driver.StaticFileSystemInterface) {
	r.APIRoute("/api/auth", NewAuthHandler)
	r.APIAuthorizedRoute("/api/content/", NewContentHandler)
	r.APIAuthorizedRoute("/api/search", NewSearchContentHandler)

	r.HandleWithCacheControl(
		"/api/uploads/",
		http.StripPrefix("/api/uploads/", http.FileServer(uploadsFileSystem)),
	)
}
