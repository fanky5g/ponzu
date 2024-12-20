package dashboard

import (
	"net/http"

	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/dashboard/content/edit"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"

	"github.com/fanky5g/ponzu/tokens"
)

func RegisterRoutes(r router.Router) error {
	contentService := r.Context().Service(tokens.ContentServiceToken).(content.Service)
	propCache := r.Context().Service(tokens.ApplicationPropertiesProviderToken).(config.ApplicationPropertiesCache)

	r.AuthorizedRoute("GET /edit", func(r router.Router) http.HandlerFunc {
		return edit.NewEditContentFormHandler(propCache, contentService)
	})

	return nil
}
