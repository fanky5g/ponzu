package api

import (
	"github.com/fanky5g/ponzu/internal/handler/controllers/middleware"
	"github.com/fanky5g/ponzu/internal/services"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/internal/services/search"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"net/http"
)

// RegisterRoutes adds Handlers to default http listener for API
func RegisterRoutes(mux *http.ServeMux, services services.Services, middlewares middleware.Middlewares) {
	// Services
	authService := services.Get(auth.ServiceToken).(auth.Service)
	contentService := services.Get(content.ServiceToken).(content.Service)
	storageService := services.Get(storage.ServiceToken).(storage.Service)
	contentSearchService := services.Get(search.ContentSearchService).(search.Service)
	// End Services

	// Middlewares
	AnalyticsRecorderMiddleware := middlewares.Get(middleware.AnalyticsRecorderMiddleware)
	CORSMiddleware := middlewares.Get(middleware.CorsMiddleware)
	GzipMiddleware := middlewares.Get(middleware.GzipMiddleware)

	Auth := middlewares.Get(middleware.AuthMiddleware)
	// End Middlewares

	mux.HandleFunc("/api/auth", AnalyticsRecorderMiddleware(CORSMiddleware(NewAuthHandler(authService))))
	mux.HandleFunc(
		"/api/content/",
		AnalyticsRecorderMiddleware(
			CORSMiddleware(Auth(GzipMiddleware(NewContentHandler(contentService, storageService)))),
		),
	)

	mux.HandleFunc("/api/search",
		AnalyticsRecorderMiddleware(
			Auth(CORSMiddleware(GzipMiddleware(NewSearchContentHandler(contentSearchService)))),
		),
	)
}
