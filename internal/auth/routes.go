package auth

import (
	"github.com/fanky5g/ponzu/internal/layouts"
	"net/http"
)

type RouterInterface interface {
	AuthorizedRoute(pattern string, handler func() http.HandlerFunc)
	APIRoute(pattern string, handler func() http.HandlerFunc)
	Route(pattern string, handler func() http.HandlerFunc)
}

func RegisterRoutes(
	r RouterInterface,
	authService *Service,
	userService *UserService,
	publicPath string,
	rootLayout layouts.Template,
	dashboardLayout layouts.Template,
) {
	r.APIRoute("POST /api/auth", func() http.HandlerFunc {
		return NewAPIAuthHandler(authService)
	})

	r.Route("GET /logout", func() http.HandlerFunc {
		return NewLogoutHandler(publicPath)
	})

	r.Route("GET /login", func() http.HandlerFunc {
		return NewLoginHandler(publicPath, authService, userService, rootLayout)
	})

	r.Route("POST /login", func() http.HandlerFunc {
		return NewAuthHandler(publicPath, authService)
	})

	r.Route("GET /recover", func() http.HandlerFunc {
		return NewForgotPasswordPageHandler(publicPath, rootLayout)
	})

	r.Route("POST /recover", func() http.HandlerFunc {
		return NewForgotPasswordHandler(publicPath, authService)
	})

	r.Route("GET /recover/key", func() http.HandlerFunc {
		return NewRecoveryKeyPageHandler(publicPath, rootLayout)
	})

	r.Route("POST /recover/key", func() http.HandlerFunc {
		return NewRecoveryKeyHandler(publicPath, authService, userService)
	})

	// Dashboard routes (Configuration)
	r.AuthorizedRoute("GET /configure/users", func() http.HandlerFunc {
		return NewUsersPageHandler(publicPath, authService, userService, dashboardLayout)
	})

	r.AuthorizedRoute("POST /configure/users", func() http.HandlerFunc {
		return NewCreateUserHandler(publicPath, authService, userService)
	})

	r.AuthorizedRoute("POST /configure/users/edit", func() http.HandlerFunc {
		return NewConfigUsersEditHandler(publicPath, authService, userService)
	})

	r.AuthorizedRoute("POST /configure/users/delete", func() http.HandlerFunc {
		return NewDeleteUserHandler(publicPath, authService, userService)
	})
}
