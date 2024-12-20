package router

import (
	"github.com/fanky5g/ponzu/internal/handler/controllers/middleware"
	"github.com/fanky5g/ponzu/internal/handler/controllers/renderer"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router/context"
	"github.com/fanky5g/ponzu/internal/application"
	"net/http"
)

type Router interface {
	Redirect(*http.Request, http.ResponseWriter, string)
	Route(string, http.HandlerFunc)
	AuthorizedRoute(string, http.HandlerFunc)
	Handle(string, http.Handler)
	HandleWithCacheControl(string, http.Handler)
	APIRoute(string, http.HandlerFunc)
	APIAuthorizedRoute(string, http.HandlerFunc)

	Context() context.Context
	Renderer() renderer.Renderer
}

type router struct {
	middlewares middleware.Middlewares
	mux         *http.ServeMux
	ctx         context.Context
	renderer    renderer.Renderer
}

func New(mux *http.ServeMux) (Router, error) {
	middlewares, err := middleware.New(
		paths,
		svcs,
	)

	if err != nil {
		return nil, err
	}

	ctx, err := context.NewContext(svcs, types, paths)
	if err != nil {
		return nil, err
	}

	rdr, err := renderer.New(ctx)
	if err != nil {
		return nil, err
	}

	return &router{
		middlewares: middlewares,
		mux:         mux,
		ctx:         ctx,
		renderer:    rdr,
	}, nil
}
