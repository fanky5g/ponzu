package setup

import (
	"github.com/fanky5g/ponzu/internal/auth"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/layouts"
	"net/http"
)

type RouterInterface interface {
	Route(pattern string, handler func() http.HandlerFunc)
}

func RegisterRoutes(
	r RouterInterface,
	configService *config.Service,
	authService *auth.Service,
	userService *auth.UserService,
	publicPath string,
	rootTemplate layouts.Template,
) {
	r.Route("GET /init", func() http.HandlerFunc {
		return NewInitPageHandler(publicPath, userService, rootTemplate)
	})

	r.Route("POST /init", func() http.HandlerFunc {
		return NewInitHandler(publicPath, configService, authService, userService)
	})
}
