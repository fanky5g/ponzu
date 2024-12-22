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

	r.AuthorizedRoute("GET /edit", func(r router.Router) http.HandlerFunc {
		return NewEditContentFormHandler(contentService, configCache, r.Context().Paths().PublicPath)
	})

	return nil
}
