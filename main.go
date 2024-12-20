package main

import (
	"github.com/fanky5g/ponzu/internal/handler/command"
)

type DatabaseConfig struct {
	Models []models.ModelInterface
}

type Config struct {
	ContentTypes content.Types
	Database     DatabaseConfig
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
	//infra, err := infrastructure.New(conf.ContentTypes.Content, conf.Database.Models)
	//if err != nil {
	//	return nil, err
	//}

	//svcs, err := services.New(infra, conf.ContentTypes.Content)
	//if err != nil {
	//	return nil, err
	//}

	//svr, err := server.New(conf.ContentTypes, infra, svcs)
	//if err != nil {
	//	return nil, err
	//}

	return &application{server: svr}, nil
}

func main() {
	command.Execute()
}
