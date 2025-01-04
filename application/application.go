package application

import (
	"fmt"

	bleveSearch "github.com/fanky5g/ponzu-driver-bleve"
	gcsStorage "github.com/fanky5g/ponzu-driver-gcs"
	localStorage "github.com/fanky5g/ponzu-driver-local-storage"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content"
	databasePkg "github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/driver"
	contentService "github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/content/dataexporter"
	"github.com/fanky5g/ponzu/internal/content/dataexporter/formatter"
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/internal/database/postgres"
	"github.com/fanky5g/ponzu/internal/http/server"
	"github.com/fanky5g/ponzu/internal/search"
	pgSearch "github.com/fanky5g/ponzu/internal/search/postgres"
	"github.com/fanky5g/ponzu/internal/services"
	"github.com/fanky5g/ponzu/internal/uploads"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrStorageDriverMissing     = errors.New("Invalid configuration: missing upload storage_driver")
	ErrUnsupportedStorageDriver = errors.New("Unsupported upload storage driver")
	ErrSearchDriverMissing      = errors.New("Invalid configuration: missing search_driver")
	ErrUnsupportedSearchDriver  = errors.New("Unsupported search driver")
)

type DatabaseConfig struct {
	Models []databasePkg.ModelInterface
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
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	db, err := getDatabaseDriver(cfg.DatabaseDriver, conf.Database.Models)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	uploadStorageClient, err := getUploadStorageClient(cfg.StorageDriver)
	if err != nil {
		return nil, err
	}

	assetStorageClient, err := localStorage.New(config.AssetStaticDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create asset storage file system: %v", err)
	}

	searchClient, err := getSearchDriver(cfg.SearchDriver, db)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get search client")
	}

	svcs, err := services.New(db, searchClient, conf.ContentTypes.Content)
	if err != nil {
		return nil, err
	}

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

	storageService, err := uploads.New(db, searchClient, uploadStorageClient)
	if err != nil {
		log.Fatalf("Failed to initialize storage services: %v", err)
	}
	svcs[tokens.UploadServiceToken] = storageService

	svr, err := server.New(conf.ContentTypes, assetStorageClient, uploadStorageClient, svcs)
	if err != nil {
		return nil, err
	}

	return &application{server: svr}, nil
}

func getUploadStorageClient(driver string) (driver.StorageClientInterface, error) {
	switch driver {
	case "":
		return nil, ErrStorageDriverMissing
	case "local":
		uploadStorageClient, err := localStorage.New("")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize storage client: %v", err)
		}

		return uploadStorageClient, nil
	case "gcs":
		gcsStorageClient, err := gcsStorage.New()
		if err != nil {
			return nil, errors.Wrap(err, "failed to initialize gcs storage driver")
		}

		return gcsStorageClient, nil
	default:
		return nil, ErrUnsupportedStorageDriver
	}
}

func getDatabaseDriver(driver string, contentModels []databasePkg.ModelInterface) (database.Database, error) {
	switch driver {
	case "postgres":
		return postgres.New(contentModels)
	default:
		return nil, errors.New("invalid driver")
	}
}

func getSearchDriver(driver string, db database.Database) (search.SearchInterface, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	if cfg.SearchDriver == "" {
		return nil, ErrSearchDriverMissing
	}

	var searchClient search.SearchInterface
	switch cfg.SearchDriver {
	case "postgres":
		return pgSearch.New(db)
	case "bleve":
		searchClient, err = bleveSearch.New(cfg.Paths.DataDir)
		if err != nil {
			return nil, errors.Wrap(err, "failed to initialize search client")
		}
	default:
		return nil, ErrUnsupportedSearchDriver
	}

	return searchClient, nil
}
