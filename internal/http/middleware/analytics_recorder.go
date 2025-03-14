package middleware

import (
	"net/http"

	"github.com/fanky5g/ponzu/internal/analytics"
	"github.com/fanky5g/ponzu/internal/http/request"
)

var AnalyticsRecorderMiddleware Token = "AnalyticsRecorderMiddleware"

type AnalyticsRecorder interface {
	Record(req analytics.HTTPRequestMetadata)
}

func NewAnalyticsRecorderMiddleware(recorder AnalyticsRecorder) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			go recorder.Record(request.GetAnalyticsRequestMetadata(req))

			next.ServeHTTP(res, req)
		}
	}
}
