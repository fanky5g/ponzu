package content

import (
	"github.com/fanky5g/ponzu/internal/dashboard"
	"net/http"

	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"

	"github.com/fanky5g/ponzu/tokens"
)

func RegisterRoutes(r router.Router, dashboardHandler dashboard.LayoutRouteHandler) error {
	uploadService := r.Context().Service(tokens.UploadServiceToken).(*UploadService)
	contentService := r.Context().Service(tokens.ContentServiceToken).(*Service)
	configCache := r.Context().Service(tokens.ConfigCache).(config.ConfigCache)
	publicPath := r.Context().Paths().PublicPath

	mapper := NewMapper()

	r.AuthorizedRoute("GET /edit", func(r router.Router) http.HandlerFunc {
		return NewEditContentFormHandler(contentService, configCache, publicPath)
	})

	r.AuthorizedRoute("POST /edit", func(r router.Router) http.HandlerFunc {
		return NewSaveContentHandler(contentService, publicPath)
	})

	r.AuthorizedRoute("POST /edit/workflow", func(r router.Router) http.HandlerFunc {
		return NewContentWorkflowTransitionHandler(contentService, configCache, publicPath)
	})

	r.AuthorizedRoute("GET /edit/upload", func(r router.Router) http.HandlerFunc {
		return dashboardHandler(NewEditUploadFormHandler(uploadService))
	})

	r.AuthorizedRoute("POST /edit/upload", func(r router.Router) http.HandlerFunc {
		return NewSaveUploadHandler(uploadService, publicPath)
	})

	r.APIAuthorizedRoute("GET /api/references", func(r router.Router) http.HandlerFunc {
		return NewListReferencesHandler(contentService, mapper)
	})

	r.APIAuthorizedRoute("GET /api/references/{id}", func(r router.Router) http.HandlerFunc {
		return NewGetReferenceHandler(contentService, mapper)
	})

	return nil
}
