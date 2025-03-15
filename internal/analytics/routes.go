package analytics

import (
	"github.com/fanky5g/ponzu/internal/layouts"
	"net/http"
)

type RouterInterface interface {
	AuthorizedRoute(pattern string, handler func() http.HandlerFunc)
}

func RegisterRoutes(r RouterInterface, analyticsService *Service, layout layouts.Template) {
	r.AuthorizedRoute("GET /{$}", func() http.HandlerFunc {
		return NewAnalyticsHandler(analyticsService, layout)
	})
}
