package api

import (
	"github.com/fanky5g/ponzu/internal/application"
	"github.com/fanky5g/ponzu/internal/application/auth"
	"github.com/fanky5g/ponzu/internal/application/config"
	"github.com/fanky5g/ponzu/internal/handler/controllers/middleware"
	"net/http"
)

// RegisterRoutes adds Handlers to default http listener for API
func RegisterRoutes(services application.Services, middlewares middleware.Middlewares) {
	authService := services.Get(auth.ServiceToken).(auth.Service)
	CacheControlMiddleware := middlewares.Get(middleware.CacheControlMiddleware)
	configService := services.Get(config.ServiceToken).(config.Service)

	AnalyticsRecorderMiddleware := middlewares.Get(middleware.AnalyticsRecorderMiddleware)
	CORS := middleware.NewCORSMiddleware(configService, CacheControlMiddleware)

	http.HandleFunc("/api/auth", AnalyticsRecorderMiddleware(CORS(NewAuthHandler(authService))))
}
