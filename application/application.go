package application

import (
	"github.com/fanky5g/ponzu/application/server"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/infrastructure"
	"github.com/fanky5g/ponzu/internal/services"
)

type Config struct {
	ServeConfig  *server.ServeConfig
	ContentTypes content.Types
}

type application struct {
	server server.Server
}

type Application interface {
	Server() server.Server
}

func (app *application) Server() server.Server {
	return app.server
}

func New(conf Config) (Application, error) {
	infra, err := infrastructure.New(conf.ContentTypes.Content)
	if err != nil {
		return nil, err
	}

	svcs, err := services.New(infra, conf.ContentTypes.Content)
	if err != nil {
		return nil, err
	}

	svr, err := server.New(conf.ServeConfig, conf.ContentTypes, infra, svcs)
	if err != nil {
		return nil, err
	}

	return &application{server: svr}, nil
}
