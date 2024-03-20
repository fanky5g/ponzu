package infrastructure

import (
	"fmt"
	bleveSearch "github.com/fanky5g/ponzu-driver-bleve"
	boltDb "github.com/fanky5g/ponzu-driver-bolt-db"
	"github.com/fanky5g/ponzu-driver-local-storage"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
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

func New(contentTypes map[string]content.Builder) (Infrastructure, error) {
	svcs := make(map[tokens.Driver]interface{})

	db, err := boltDb.New(contentTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	uploadStorageClient, err := storage.New(config.UploadDir())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage client: %v", err)
	}

	assetStorageClient, err := storage.New(config.AssetStaticDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create asset storage file system: %v", err)
	}

	contentSearchClient, err := bleveSearch.New(
		db.Get(tokens.ContentRepositoryToken).(repositories.ContentRepositoryInterface),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize search client: %v", err)
	}

	uploadsSearchClient, err := bleveSearch.New(
		db.Get(tokens.UploadRepositoryToken).(repositories.ContentRepositoryInterface),
	)
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
