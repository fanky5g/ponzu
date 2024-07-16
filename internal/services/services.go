package services

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/infrastructure"
	"github.com/fanky5g/ponzu/internal/services/analytics"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/services/config"
	contentService "github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/internal/services/search"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/internal/services/tls"
	"github.com/fanky5g/ponzu/internal/services/users"
	"github.com/fanky5g/ponzu/tokens"
	"log"
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
	infra infrastructure.Infrastructure,
	types map[string]content.Builder,
) (Services, error) {
	db := infra.Get(tokens.DatabaseInfrastructureToken).(driver.Database)
	searchClient := infra.Get(tokens.SearchClientInfrastructureToken).(driver.SearchInterface)
	storageClient := infra.Get(tokens.StorageClientInfrastructureToken).(driver.StorageClientInterface)

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

	configService, err := config.New(db)
	if err != nil {
		log.Fatalf("Failed to initialize config services: %v", err)
	}
	services[tokens.ConfigServiceToken] = configService

	contentSvc, err := contentService.New(db, types, searchClient)
	if err != nil {
		log.Fatalf("Failed to initialize entities service: %v", err)
	}
	services[tokens.ContentServiceToken] = contentSvc

	contentSearchService, err := search.New(searchClient, db)
	if err != nil {
		log.Fatalf("Failed to initialize search service: %v", err)
	}
	services[tokens.ContentSearchServiceToken] = contentSearchService

	uploadSearchService, err := search.New(searchClient, db)
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
