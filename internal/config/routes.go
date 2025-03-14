package config

import (
	"github.com/fanky5g/ponzu/internal/layouts"
	"net/http"
)

type RouterInterface interface {
	AuthorizedRoute(pattern string, handler func() http.HandlerFunc)
}

func RegisterRoutes(r RouterInterface, publicPath string, configService *Service, layout layouts.Template) {
	r.AuthorizedRoute("GET /configure", func() http.HandlerFunc {
		return NewEditConfigHandler(publicPath, configService, layout)
	})

	r.AuthorizedRoute("POST /configure", func() http.HandlerFunc {
		return NewSaveConfigHandler(configService, publicPath)
	})
}
