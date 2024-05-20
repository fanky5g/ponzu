package infrastructure

import (
	"errors"
	"fmt"
	bleveSearch "github.com/fanky5g/ponzu-driver-bleve"
	"github.com/fanky5g/ponzu-driver-local-storage"
	postgres "github.com/fanky5g/ponzu-driver-postgres/database"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/models"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

type (
	Infrastructure interface {
		Get(token tokens.Driver) interface{}
	}
	infrastructure struct {
		services map[tokens.Driver]interface{}
	}
)

func (infra *infrastructure) Get(token tokens.Driver) interface{} {
	if service, ok := infra.services[token]; ok {
		return service
	}

	log.Fatalf("Service %s is not implemented", token)
	return nil
}

func getDatabaseDriver(
	name string,
	contentModels []models.ModelInterface,
) (driver.Database, error) {
	m := make([]models.ModelInterface, 0)
	m = append(m, contentModels...)
	m = append(m, models.GetPonzuModels()...)

	switch name {
	case "postgres":
		return postgres.New(m)
	default:
		return nil, errors.New("invalid driver")
	}
}

func New(
	contentTypes map[string]content.Builder,
	contentModels []models.ModelInterface,
) (Infrastructure, error) {
	svcs := make(map[tokens.Driver]interface{})
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	var db driver.Database
	db, err = getDatabaseDriver(cfg.DatabaseDriver, contentModels)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	uploadStorageClient, err := storage.New("")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage client: %v", err)
	}

	assetStorageClient, err := storage.New(config.AssetStaticDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create asset storage file system: %v", err)
	}

	contentSearchClient, err := bleveSearch.New(
		contentTypes,
		db,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize search client: %v", err)
	}

	uploadsSearchClient, err := bleveSearch.New(map[string]content.Builder{
		constants.UploadsEntityName: func() interface{} {
			return new(entities.FileUpload)
		},
	}, db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize upload search client: %v", err)
	}

	svcs[tokens.StorageClientInfrastructureToken] = uploadStorageClient
	svcs[tokens.AssetStorageClientInfrastructureToken] = assetStorageClient
	svcs[tokens.ContentSearchClientInfrastructureToken] = contentSearchClient
	svcs[tokens.UploadSearchClientInfrastructureToken] = uploadsSearchClient
	svcs[tokens.DatabaseInfrastructureToken] = db

	return &infrastructure{services: svcs}, nil
}
