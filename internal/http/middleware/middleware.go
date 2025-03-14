package middleware

import (
	conf "github.com/fanky5g/ponzu/config"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ConfigProviderInterface interface {
	CacheControlConfigInterface
	CorsConfigInterface
	GzipConfigInterface
}

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

func New(
	paths conf.Paths,
	tokenValidator TokenValidatorInterface,
	configProvider ConfigProviderInterface,
	analyticsRecorder AnalyticsRecorder,
) (Middlewares, error) {
	middlewares := make(Middlewares)

	cacheControlMiddleware := NewCacheControlMiddleware(configProvider)
	middlewares[CacheControlMiddleware] = cacheControlMiddleware
	middlewares[AnalyticsRecorderMiddleware] = NewAnalyticsRecorderMiddleware(analyticsRecorder)
	middlewares[AuthMiddleware] = NewAuthMiddleware(paths, tokenValidator)
	middlewares[GzipMiddleware] = NewGzipMiddleware(configProvider)
	middlewares[CorsMiddleware] = NewCORSMiddleware(configProvider, cacheControlMiddleware)

	return middlewares, nil
}
