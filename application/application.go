package application

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/infrastructure"
	contentService "github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/content/dataexporter"
	"github.com/fanky5g/ponzu/internal/content/dataexporter/formatter"
	"github.com/fanky5g/ponzu/internal/http/server"
	"github.com/fanky5g/ponzu/internal/services"
	"github.com/fanky5g/ponzu/models"
	"github.com/fanky5g/ponzu/tokens"

	log "github.com/sirupsen/logrus"
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
	infra, err := infrastructure.New(conf.ContentTypes.Content, conf.Database.Models)
	if err != nil {
		return nil, err
	}

	svcs, err := services.New(infra, conf.ContentTypes.Content)
	if err != nil {
		return nil, err
	}

	db := infra.Get(tokens.DatabaseInfrastructureToken).(driver.Database)
	searchClient := infra.Get(tokens.SearchClientInfrastructureToken).(driver.SearchInterface)

	rowFormatter, err := formatter.NewJSONRowFormatter()
	if err != nil {
		return nil, err
	}

	contentExporter, err := dataexporter.New(rowFormatter)
	if err != nil {
		return nil, err
	}

	contentSvc, err := contentService.New(
		db,
		conf.ContentTypes.Content,
		searchClient,
		contentExporter,
	)
	if err != nil {
		log.Fatalf("Failed to initialize entities service: %v", err)
	}
	svcs[tokens.ContentServiceToken] = contentSvc

	svr, err := server.New(conf.ContentTypes, infra, svcs)
	if err != nil {
		return nil, err
	}

	return &application{server: svr}, nil
}
