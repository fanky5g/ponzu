package middleware

import (
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/services"
	"github.com/fanky5g/ponzu/internal/services/analytics"
	"github.com/fanky5g/ponzu/tokens"
	"log"
	"net/http"
)

type Token string

type Middleware func(next http.HandlerFunc) http.HandlerFunc
type Middlewares map[Token]Middleware

func (middlewares Middlewares) Get(token Token) Middleware {
	if middleware, ok := middlewares[token]; ok {
		return middleware
	}

	log.Fatalf("Middleware %s is not implemented", token)
	return nil
}

func New(paths conf.Paths, applicationServices services.Services) (Middlewares, error) {
	middlewares := make(Middlewares)
	analyticsService := applicationServices.Get(tokens.AnalyticsServiceToken).(analytics.Service)
	configCache := applicationServices.Get(tokens.ConfigCache).(config.ConfigCache)

	cacheControlMiddleware := NewCacheControlMiddleware(configCache)
	middlewares[CacheControlMiddleware] = cacheControlMiddleware
	middlewares[AnalyticsRecorderMiddleware] = NewAnalyticsRecorderMiddleware(analyticsService)
	middlewares[AuthMiddleware] = NewAuthMiddleware(paths, applicationServices)
	middlewares[GzipMiddleware] = NewGzipMiddleware(configCache)
	middlewares[CorsMiddleware] = NewCORSMiddleware(configCache, cacheControlMiddleware)

	return middlewares, nil
}
