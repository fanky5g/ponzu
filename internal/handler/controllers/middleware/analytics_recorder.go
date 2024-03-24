package middleware

import (
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/services/analytics"
	"net/http"
)

var AnalyticsRecorderMiddleware Token = "AnalyticsRecorderMiddleware"

func NewAnalyticsRecorderMiddleware(analyticsService analytics.Service) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			go analyticsService.Record(request.GetAnalyticsRequestMetadata(req))

			next.ServeHTTP(res, req)
		}
	}
}
