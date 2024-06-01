package middleware

import (
	"fmt"
	"github.com/fanky5g/ponzu/internal/services/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	// DefaultMaxAge provides a 2592000 second (30-day) cache max-age setting
	DefaultMaxAge = int64(60 * 60 * 24 * 30)
)

var CacheControlMiddleware Token = "CacheControlMiddleware"

func NewCacheControlMiddleware(configService config.Service) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {

		return func(res http.ResponseWriter, req *http.Request) {
			cfg, err := configService.Get()
			if err != nil {
				log.WithField("Error", err).Warning("Failed to get get config")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			if cfg.DisableHTTPCache {
				res.Header().Add("Cache-Control", "no-cache")
				next.ServeHTTP(res, req)
			} else {
				age := cfg.CacheMaxAge
				etag := cfg.Etag

				if age == 0 {
					age = DefaultMaxAge
				}

				policy := fmt.Sprintf("max-age=%d, public", age)
				res.Header().Add("ETag", etag)
				res.Header().Add("Cache-Control", policy)

				if match := req.Header.Get("If-None-Match"); match != "" {
					if strings.Contains(match, etag) {
						res.WriteHeader(http.StatusNotModified)
						return
					}
				}

				next.ServeHTTP(res, req)
			}
		}
	}
}

func ToHttpHandler(middleware Middleware) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return middleware(func(res http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(res, req)
		})
	}
}
