package uploads

import (
	"net/http"

	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/tokens"
)

func RegisterRoutes(r router.Router) {
	uploadService := r.Context().Service(tokens.UploadServiceToken).(*Service)
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)
	configCache := r.Context().Service(tokens.ConfigCache).(config.ConfigCache)
	publicPath := r.Context().Paths().PublicPath

	r.AuthorizedRoute("GET /edit/upload", func(r router.Router) http.HandlerFunc {
		return NewEditUploadFormHandler(uploadService, contentService, configCache, publicPath)
	})

	r.AuthorizedRoute("POST /edit/upload", func(r router.Router) http.HandlerFunc {
		return NewSaveUploadHandler(uploadService, publicPath)
	})
}
