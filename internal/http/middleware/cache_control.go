package middleware

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	// DefaultMaxAge provides a 2592000 second (30-day) cache max-age setting
	DefaultMaxAge = int64(60 * 60 * 24 * 30)
)

type CacheControlConfigInterface interface {
	GetHTTPCacheDisabled() (bool, error)
	GetCacheControlMaxAge() (int64, error)
}

var CacheControlMiddleware Token = "CacheControlMiddleware"

func NewCacheControlMiddleware(cacheControlConfig CacheControlConfigInterface) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			httpCacheDisabled, err := cacheControlConfig.GetHTTPCacheDisabled()
			if err != nil {
				log.WithField("Error", err).Warning("Failed to get config")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			if httpCacheDisabled {
				res.Header().Add("Cache-Control", "no-cache")
				next.ServeHTTP(res, req)
			} else {
				cacheMaxAge, err := cacheControlConfig.GetCacheControlMaxAge()
				if err != nil {
					log.WithField("Error", err).Warning("Failed to get config")
					res.WriteHeader(http.StatusInternalServerError)
					return
				}

				if cacheMaxAge == 0 {
					cacheMaxAge = DefaultMaxAge
				}

				policy := fmt.Sprintf("max-age=%d, public", cacheMaxAge)
				res.Header().Add("Cache-Control", policy)

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
