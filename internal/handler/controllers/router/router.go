package router

import (
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/middleware"
	"github.com/fanky5g/ponzu/internal/handler/controllers/renderer"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router/context"
	"github.com/fanky5g/ponzu/internal/services"
	"net/http"
)

type Router interface {
	Redirect(*http.Request, http.ResponseWriter, string)
	Route(string, func(r Router) http.HandlerFunc)
	AuthorizedRoute(string, func(r Router) http.HandlerFunc)
	Handle(pattern string, handler http.Handler)
	HandleWithCacheControl(pattern string, handler http.Handler)
	APIRoute(string, func(r Router) http.HandlerFunc)
	APIAuthorizedRoute(pattern string, handler func(r Router) http.HandlerFunc)

	Context() context.Context
	Renderer() renderer.Renderer
}

type router struct {
	middlewares middleware.Middlewares
	mux         *http.ServeMux
	ctx         context.Context
	renderer    renderer.Renderer
}

func (r *router) Context() context.Context {
	return r.ctx
}

func (r *router) Renderer() renderer.Renderer {
	return r.renderer
}

func New(
	mux *http.ServeMux,
	paths conf.Paths,
	svcs services.Services,
	types content.Types) (Router, error) {
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
