package api

import (
	"github.com/fanky5g/ponzu/internal/application"
	"github.com/fanky5g/ponzu/internal/application/auth"
	"github.com/fanky5g/ponzu/internal/application/config"
	"github.com/fanky5g/ponzu/internal/application/content"
	"github.com/fanky5g/ponzu/internal/application/storage"
	"github.com/fanky5g/ponzu/internal/handler/controllers/middleware"
	"net/http"
)

// RegisterRoutes adds Handlers to default http listener for API
func RegisterRoutes(services application.Services, middlewares middleware.Middlewares) {
	// Services
	authService := services.Get(auth.ServiceToken).(auth.Service)
	configService := services.Get(config.ServiceToken).(config.Service)
	contentService := services.Get(content.ServiceToken).(content.Service)
	storageService := services.Get(storage.ServiceToken).(storage.Service)
	// End Services

	// Middlewares
	CacheControlMiddleware := middlewares.Get(middleware.CacheControlMiddleware)
	AnalyticsRecorderMiddleware := middlewares.Get(middleware.AnalyticsRecorderMiddleware)
	CORS := middleware.NewCORSMiddleware(configService, CacheControlMiddleware)
	Auth := middlewares.Get(middleware.AuthMiddleware)
	// End Middlewares

	http.HandleFunc("/api/auth", AnalyticsRecorderMiddleware(CORS(NewAuthHandler(authService))))
	http.HandleFunc(
		"/api/content",
		AnalyticsRecorderMiddleware(CORS(Auth(NewCreateContentHandler(contentService, storageService)))),
	)
}
