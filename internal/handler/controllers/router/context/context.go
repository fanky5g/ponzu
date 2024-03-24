package context

import (
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/services"
	"github.com/fanky5g/ponzu/tokens"
)

type context struct {
	services services.Services
	types    content.Types
	paths    config.Paths
}

type Context interface {
	Service(name tokens.Service) interface{}
	Types() content.Types
	Paths() config.Paths
}

func (ctx *context) Service(name tokens.Service) interface{} {
	return ctx.services.Get(name)
}

func (ctx *context) Types() content.Types {
	return ctx.types
}

func (ctx *context) Paths() config.Paths {
	return ctx.paths
}

func NewContext(
	s services.Services,
	types content.Types,
	paths config.Paths,
) (Context, error) {
	return &context{
		services: s,
		types:    types,
		paths:    paths,
	}, nil
}
