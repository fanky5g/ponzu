package infrastructure

import (
	"fmt"

    "github.com/pkg/errors"
	localStorage "github.com/fanky5g/ponzu-driver-local-storage"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/driver"
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

    uploadStorageClient, err := getUploadStorageClient()
    if err != nil {
        return nil, err
    }

	assetStorageClient, err := localStorage.New(config.AssetStaticDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create asset storage file system: %v", err)
	}

    searchClient, err := getSearchClient(db)
    if err != nil {
        return nil, errors.Wrap(err, "Failed to get search client")
    }

	svcs[tokens.StorageClientInfrastructureToken] = uploadStorageClient
	svcs[tokens.AssetStorageClientInfrastructureToken] = assetStorageClient
	svcs[tokens.SearchClientInfrastructureToken] = searchClient
	svcs[tokens.DatabaseInfrastructureToken] = db

	return &infrastructure{services: svcs}, nil
}
