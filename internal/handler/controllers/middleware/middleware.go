package middleware

import (
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
	"github.com/fanky5g/ponzu/internal/services"
	"github.com/fanky5g/ponzu/internal/services/analytics"
	"github.com/fanky5g/ponzu/internal/services/config"
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

func New(paths conf.Paths, applicationServices services.Services, cache repositories.Cache) (Middlewares, error) {
	middlewares := make(Middlewares)
	analyticsService := applicationServices.Get(tokens.AnalyticsServiceToken).(analytics.Service)
	configService := applicationServices.Get(tokens.ConfigServiceToken).(config.Service)

	cacheControlMiddleware := NewCacheControlMiddleware(cache)
	middlewares[CacheControlMiddleware] = cacheControlMiddleware
	middlewares[AnalyticsRecorderMiddleware] = NewAnalyticsRecorderMiddleware(analyticsService)
	middlewares[AuthMiddleware] = NewAuthMiddleware(paths, applicationServices)
	middlewares[GzipMiddleware] = NewGzipMiddleware(configService)
	middlewares[CorsMiddleware] = NewCORSMiddleware(applicationServices, cacheControlMiddleware)

	return middlewares, nil
}
