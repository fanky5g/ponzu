package router

import (
	"github.com/fanky5g/ponzu/internal/http/middleware"
	"net/http"
)

type Router struct {
	middlewares middleware.Middlewares
	mux         *http.ServeMux
}

func New(mux *http.ServeMux, middlewares middleware.Middlewares) (*Router, error) {
	return &Router{
		middlewares: middlewares,
		mux:         mux,
	}, nil
}
