package router

import (
	"context"
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/fanky5g/ponzu/internal/http/middleware"
	"net/http"
)

func (r *Router) Route(pattern string, handler func() http.HandlerFunc) {
	r.mux.HandleFunc(pattern, handler())
}

func (r *Router) APIRoute(pattern string, handler func() http.HandlerFunc) {
	AnalyticsRecorderMiddleware := r.middlewares.Get(middleware.AnalyticsRecorderMiddleware)
	CORSMiddleware := r.middlewares.Get(middleware.CorsMiddleware)
	GzipMiddleware := r.middlewares.Get(middleware.GzipMiddleware)
	TagRoute := r.RouteTag(constants.APIRoute)

	r.mux.HandleFunc(pattern, TagRoute(AnalyticsRecorderMiddleware(CORSMiddleware(GzipMiddleware(handler())))))
}

func (r *Router) APIAuthorizedRoute(pattern string, handler func() http.HandlerFunc) {
	Auth := r.middlewares.Get(middleware.AuthMiddleware)

	r.APIRoute(pattern, func() http.HandlerFunc {
		return Auth(handler())
	})
}

func (r *Router) AuthorizedRoute(pattern string, handler func() http.HandlerFunc) {
	Auth := r.middlewares.Get(middleware.AuthMiddleware)
	r.mux.HandleFunc(pattern, Auth(handler()))
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

func (r *Router) HandleWithCacheControl(pattern string, handler http.Handler) {
	CacheControlMiddleware := middleware.ToHttpHandler(r.middlewares.Get(middleware.CacheControlMiddleware))
	r.Handle(pattern, CacheControlMiddleware(handler))
}

func (r *Router) RouteTag(value constants.RouteTag) func(next http.HandlerFunc) http.HandlerFunc {
	return r.Tag(constants.RouteTagIdentifier, value)
}

func (r *Router) Tag(key string, value any) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(
				res,
				req.WithContext(context.WithValue(req.Context(), key, value)),
			)
		}
	}
}
