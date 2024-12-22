package middleware

import (
	"github.com/fanky5g/ponzu/internal/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

var CorsMiddleware Token = "CorsMiddleware"

// sendPreflight is used to respond to a cross-origin "OPTIONS" request
func sendPreflight(res http.ResponseWriter) {
	res.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.WriteHeader(200)
	return
}

func responseWithCORS(
	corsDisabled bool,
	domain string,
	res http.ResponseWriter,
	req *http.Request) (http.ResponseWriter, bool) {
	if corsDisabled {
		origin := req.Header.Get("Origin")
		u, err := url.Parse(origin)
		if err != nil {
			log.Println("Error parsing URL from request Origin header:", origin)
			return res, false
		}

		// hack to get dev environments to bypass cors since u.Host (below) will
		// be empty, based on Go's url.Parse function
		if domain == "localhost" {
			domain = ""
		}
		origin = u.Host

		// currently, this will check for exact match. will need feedback to
		// determine if subdomains should be allowed or allow multiple domains
		// in config
		if origin == domain {
			// apply limited CORS headers and return
			res.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
			res.Header().Set("Access-Control-Allow-Origin", domain)
			return res, true
		}

		// disallow request
		res.WriteHeader(http.StatusForbidden)
		return res, false
	}

	// apply full CORS headers and return
	res.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	return res, true
}

func NewCORSMiddleware(propCache config.ConfigCache, cacheControlMiddleware Middleware) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return cacheControlMiddleware(
			func(res http.ResponseWriter, req *http.Request) {
				corsDisabled, err := propCache.GetCorsDisabled()
				if err != nil {
					log.WithField("Error", err).Warning("Failed to get get config")
					res.WriteHeader(http.StatusInternalServerError)
					return
				}

				domain, err := propCache.GetDomain()
				if err != nil {
					log.WithField("Error", err).Warning("Failed to get get config")
					res.WriteHeader(http.StatusInternalServerError)
					return
				}

				res, cors := responseWithCORS(corsDisabled, domain, res, req)
				if !cors {
					return
				}

				if req.Method == http.MethodOptions {
					sendPreflight(res)
					return
				}

				next.ServeHTTP(res, req)
			},
		)
	}
}
