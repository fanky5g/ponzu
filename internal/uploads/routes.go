package uploads

import (
	"github.com/fanky5g/ponzu/internal/dashboard"
	"net/http"

	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/tokens"
)

func RegisterRoutes(r router.Router, dashboardHandler dashboard.LayoutRouteHandler) {
	uploadService := r.Context().Service(tokens.UploadServiceToken).(*UploadService)
	publicPath := r.Context().Paths().PublicPath

	r.AuthorizedRoute("GET /edit/upload", func(r router.Router) http.HandlerFunc {
		return dashboardHandler(NewEditUploadFormHandler(uploadService))
	})

	r.AuthorizedRoute("POST /edit/upload", func(r router.Router) http.HandlerFunc {
		return NewSaveUploadHandler(uploadService, publicPath)
	})
}
