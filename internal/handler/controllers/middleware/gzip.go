package middleware

import (
	"compress/gzip"
	"github.com/fanky5g/ponzu/internal/services/config"
	"log"
	"net/http"
	"strings"
)

var GzipMiddleware Token = "GzipMiddleware"

type gzipResponseWriter struct {
	http.ResponseWriter
	pusher http.Pusher

	gw *gzip.Writer
}

func (gzw gzipResponseWriter) Write(p []byte) (int, error) {
	return gzw.gw.Write(p)
}

func (gzw gzipResponseWriter) Push(target string, opts *http.PushOptions) error {
	if gzw.pusher == nil {
		return nil
	}

	if opts == nil {
		opts = &http.PushOptions{
			Header: make(http.Header),
		}
	}

	opts.Header.Set("Accept-Encoding", "gzip")
	return gzw.pusher.Push(target, opts)
}

func NewGzipMiddleware(configService config.Service) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			gzipDisabled, err := configService.GetCacheBoolValue("gzip_disabled")
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				log.Printf("Failed to get cache value gzip_disabled: %v\n", err)
				return
			}

			if gzipDisabled {
				next.ServeHTTP(res, req)
				return
			}

			// check if req header entities-encoding supports gzip
			if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
				// gzip response data
				res.Header().Set("Content-Encoding", "gzip")
				gzWriter := gzip.NewWriter(res)
				defer func(gzWriter *gzip.Writer) {
					err = gzWriter.Close()
					if err != nil {
						log.Printf("Failed to close gzip writer: %v\n", err)
					}
				}(gzWriter)

				var gzres gzipResponseWriter
				if pusher, ok := res.(http.Pusher); ok {
					gzres = gzipResponseWriter{res, pusher, gzWriter}
				} else {
					gzres = gzipResponseWriter{res, nil, gzWriter}
				}

				next.ServeHTTP(gzres, req)
				return
			}

			next.ServeHTTP(res, req)
		}
	}
}
