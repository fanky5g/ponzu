package services

import (
	"log"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/internal/memorycache"
	"github.com/fanky5g/ponzu/internal/search"
	"github.com/fanky5g/ponzu/internal/services/analytics"
	"github.com/fanky5g/ponzu/internal/services/auth"
	searchSvc "github.com/fanky5g/ponzu/internal/services/search"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/internal/services/tls"
	"github.com/fanky5g/ponzu/internal/services/users"
	"github.com/fanky5g/ponzu/tokens"
)

type Services map[tokens.Service]interface{}

func (services Services) Get(token tokens.Service) interface{} {
	if service, ok := services[token]; ok {
		return service
	}

	log.Fatalf("Service %s is not implemented", token)
	return nil
}

func New(
	db database.Database,
	searchClient search.SearchInterface,
	storageClient driver.StorageClientInterface,
	types map[string]content.Builder,
) (Services, error) {
	// Initialize services
	services := make(Services)

	tlsService, err := tls.New(db)
	if err != nil {
		log.Fatalf("Failed to initialize tls services %v", err)
	}
	services[tokens.TLSServiceToken] = tlsService

	userService, err := users.New(db)
	if err != nil {
		log.Fatalf("Failed to initialize user services: %v", err)
	}
	services[tokens.UserServiceToken] = userService

	authService, err := auth.New(db)

	if err != nil {
		log.Fatalf("Failed to initialize auth services: %v", err)
	}
	services[tokens.AuthServiceToken] = authService

	analyticsService, err := analytics.New(db)
	if err != nil {
		log.Fatalf("Failed to initialize analytics services: %v", err)
	}
	services[tokens.AnalyticsServiceToken] = analyticsService

	memcache, err := memorycache.New()
	if err != nil {
		log.Fatalf("Failed to initialize memory cache: %v", err)
	}

	configRepository := db.GetRepositoryByToken(tokens.ConfigRepositoryToken)
	configService, err := config.New(configRepository, memcache)
	if err != nil {
		log.Fatalf("Failed to initialize config services: %v", err)
	}
	services[tokens.ConfigServiceToken] = configService

	configCache, err := config.NewCache(memcache, types)
	if err != nil {
		log.Fatalf("Failed to initialize config cache: %v", err)
	}
	services[tokens.ConfigCache] = configCache

	contentSearchService, err := searchSvc.New(searchClient, db)
	if err != nil {
		log.Fatalf("Failed to initialize search service: %v", err)
	}
	services[tokens.ContentSearchServiceToken] = contentSearchService

	uploadSearchService, err := searchSvc.New(searchClient, db)
	if err != nil {
		log.Fatalf("Failed to initialize search service: %v", err)
	}
	services[tokens.UploadSearchServiceToken] = uploadSearchService

	storageService, err := storage.New(db, searchClient, storageClient)
	if err != nil {
		log.Fatalf("Failed to initialize storage services: %v", err)
	}
	services[tokens.StorageServiceToken] = storageService

	return services, nil
}
