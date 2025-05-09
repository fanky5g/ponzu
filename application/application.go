package application

import (
	"fmt"
	bleveSearch "github.com/fanky5g/ponzu-driver-bleve"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/workflow"
	databasePkg "github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/internal/database/postgres"
	"github.com/fanky5g/ponzu/internal/http/server"
	"github.com/fanky5g/ponzu/internal/search"
	pgSearch "github.com/fanky5g/ponzu/internal/search/postgres"
	"github.com/fanky5g/ponzu/internal/storage"
	"github.com/fanky5g/ponzu/internal/storage/assets"
	"github.com/fanky5g/ponzu/internal/storage/gcs"
	"github.com/fanky5g/ponzu/internal/storage/localstorage"
	"github.com/pkg/errors"
	"net/http"
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
	ContentTypes               content.Types
	Database                   DatabaseConfig
	ServeMux                   *http.ServeMux
	WorkflowStateChangeHandler workflow.StateChangeTrigger
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

	searchClient, err := getSearchDriver(db)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get search client")
	}

	svr, err := server.New(
		conf.ContentTypes.Content,
		db,
		assets.AssetStorage,
		uploadStorageClient,
		searchClient,
		conf.ServeMux,
		conf.WorkflowStateChangeHandler,
	)

	if err != nil {
		return nil, err
	}

	return &application{server: svr}, nil
}

func getUploadStorageClient(driver string) (storage.Client, error) {
	switch driver {
	case "":
		return nil, ErrStorageDriverMissing
	case "local":
		uploadStorageClient, err := localstorage.New("")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize storage client: %v", err)
		}

		return uploadStorageClient, nil
	case "gcs":
		gcsStorageClient, err := gcs.New()
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

func getSearchDriver(db database.Database) (search.SearchInterface, error) {
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
