package content

import (
	"net/http"

	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"

	"github.com/fanky5g/ponzu/tokens"
)

func RegisterRoutes(r router.Router) error {
	contentService := r.Context().Service(tokens.ContentServiceToken).(*Service)
	configCache := r.Context().Service(tokens.ConfigCache).(config.ConfigCache)
	publicPath := r.Context().Paths().PublicPath

	r.AuthorizedRoute("GET /edit", func(r router.Router) http.HandlerFunc {
		return NewEditContentFormHandler(contentService, configCache, publicPath)
	})

	r.AuthorizedRoute("POST /edit", func(r router.Router) http.HandlerFunc {
		return NewSaveContentHandler(contentService, publicPath)
	})

	r.AuthorizedRoute("POST /edit/workflow", func(r router.Router) http.HandlerFunc {
		return NewContentWorkflowTransitionHandler(contentService, publicPath)
	})

	return nil
}
