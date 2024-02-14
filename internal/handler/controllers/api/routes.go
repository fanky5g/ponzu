package api

import (
	"github.com/fanky5g/ponzu/internal/application"
	"github.com/fanky5g/ponzu/internal/application/auth"
	"github.com/fanky5g/ponzu/internal/application/content"
	"github.com/fanky5g/ponzu/internal/application/search"
	"github.com/fanky5g/ponzu/internal/application/storage"
	"github.com/fanky5g/ponzu/internal/handler/controllers/middleware"
	"net/http"
)

// RegisterRoutes adds Handlers to default http listener for API
func RegisterRoutes(services application.Services, middlewares middleware.Middlewares) {
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

	http.HandleFunc("/api/auth", AnalyticsRecorderMiddleware(CORSMiddleware(NewAuthHandler(authService))))
	http.HandleFunc(
		"/api/content/",
		AnalyticsRecorderMiddleware(
			CORSMiddleware(Auth(GzipMiddleware(NewContentHandler(contentService, storageService)))),
		),
	)

	http.HandleFunc("/api/search",
		AnalyticsRecorderMiddleware(
			Auth(CORSMiddleware(GzipMiddleware(NewSearchContentHandler(contentSearchService)))),
		),
	)
}
