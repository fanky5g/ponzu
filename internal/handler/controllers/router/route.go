package router

import (
	"context"
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/internal/handler/controllers/middleware"
	"net/http"
)

func (r *router) Route(pattern string, handler func(r Router) http.HandlerFunc) {
	r.mux.HandleFunc(pattern, handler(r))
}

func (r *router) APIRoute(pattern string, handler func(r Router) http.HandlerFunc) {
	AnalyticsRecorderMiddleware := r.middlewares.Get(middleware.AnalyticsRecorderMiddleware)
	CORSMiddleware := r.middlewares.Get(middleware.CorsMiddleware)
	GzipMiddleware := r.middlewares.Get(middleware.GzipMiddleware)
	TagRoute := r.RouteTag(constants.APIRoute)

	r.mux.HandleFunc(pattern, TagRoute(AnalyticsRecorderMiddleware(CORSMiddleware(GzipMiddleware(handler(r))))))
}

func (r *router) APIAuthorizedRoute(pattern string, handler func(r Router) http.HandlerFunc) {
	Auth := r.middlewares.Get(middleware.AuthMiddleware)

	r.APIRoute(pattern, func(r Router) http.HandlerFunc {
		return Auth(handler(r))
	})
}

func (r *router) AuthorizedRoute(pattern string, handler func(r Router) http.HandlerFunc) {
	Auth := r.middlewares.Get(middleware.AuthMiddleware)
	r.mux.HandleFunc(pattern, Auth(handler(r)))
}

func (r *router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

func (r *router) HandleWithCacheControl(pattern string, handler http.Handler) {
	CacheControlMiddleware := middleware.ToHttpHandler(r.middlewares.Get(middleware.CacheControlMiddleware))
	r.Handle(pattern, CacheControlMiddleware(handler))
}

func (r *router) RouteTag(value constants.RouteTag) func(next http.HandlerFunc) http.HandlerFunc {
	return r.Tag(constants.RouteTagIdentifier, value)
}

func (r *router) Tag(key string, value any) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(
				res,
				req.WithContext(context.WithValue(req.Context(), key, value)),
			)
		}
	}
}
