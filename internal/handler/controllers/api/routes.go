package api

import (
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
)

func RegisterRoutes(r router.Router) {
	r.APIRoute("/api/auth", NewAuthHandler)
	r.APIAuthorizedRoute("/api/content/", NewContentHandler)
	r.APIAuthorizedRoute("/api/search", NewSearchContentHandler)
}
