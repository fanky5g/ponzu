package application

import (
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/database"
	bleveSearchClient "github.com/fanky5g/ponzu/driver/search"
	localStorageClient "github.com/fanky5g/ponzu/driver/storage"
	"github.com/fanky5g/ponzu/internal/handler/controllers"
	"github.com/fanky5g/ponzu/internal/handler/controllers/api"
	"github.com/fanky5g/ponzu/internal/handler/controllers/middleware"
	"github.com/fanky5g/ponzu/internal/services"
	"github.com/fanky5g/ponzu/internal/services/analytics"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/services/config"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/internal/services/search"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/internal/services/tls"
	"github.com/fanky5g/ponzu/internal/services/users"
	"log"
	"net/http"
)

type ServeConfig struct {
	HttpsPort int
	HttpPort  int
	Bind      string
	DevHttps  bool
	Https     bool
}

type Config struct {
	ServeConfig *ServeConfig
	ServeMux    *http.ServeMux
	Database    database.Database
	Paths       conf.Paths
}

type application struct {
	mux *http.ServeMux

	repositories *database.Repositories
	services     services.Services
	middlewares  middleware.Middlewares
	serveConfig  *ServeConfig
}

type Application interface {
	Serve() error
	ServeMux() *http.ServeMux
}

func New(conf Config) (Application, error) {
	mux := http.NewServeMux()

	repositories, err := conf.Database.GetRepositories()
	if err != nil {
		return nil, err
	}

	// Initialize clients
	storageClient, err := localStorageClient.New()
	if err != nil {
		log.Fatalf("Failed to initialize storage client: %v", err)
	}

	contentSearchClient, err := bleveSearchClient.New(repositories.Content)
	if err != nil {
		log.Fatalf("Failed to initialize search client: %v\n", err)
	}

	uploadsSearchClient, err := bleveSearchClient.New(repositories.Uploads)
	if err != nil {
		log.Fatalf("Failed to initialize upload search client")
	}
	// End initialize clients

	// Initialize applicationServices
	applicationServices := make(services.Services)

	tlsService, err := tls.New(repositories.Config)
	if err != nil {
		log.Fatalf("Failed to initialize tls applicationServices %v", err)
	}
	applicationServices[tls.ServiceToken] = tlsService

	userService, err := users.New(repositories.Users)
	if err != nil {
		log.Fatalf("Failed to initialize user applicationServices: %v", err)
	}
	applicationServices[users.ServiceToken] = userService

	authService, err := auth.New(
		repositories.Config,
		repositories.Users,
		repositories.CredentialHashes,
		repositories.RecoveryKeys,
	)

	if err != nil {
		log.Fatalf("Failed to initialize auth applicationServices: %v", err)
	}
	applicationServices[auth.ServiceToken] = authService

	analyticsService, err := analytics.New(repositories.Analytics)
	if err != nil {
		log.Fatalf("Failed to initialize analytics applicationServices: %v", err)
	}
	applicationServices[analytics.ServiceToken] = analyticsService

	configService, err := config.New(repositories.Config)
	if err != nil {
		log.Fatalf("Failed to initialize config applicationServices: %v", err)
	}
	applicationServices[config.ServiceToken] = configService

	contentService, err := content.New(repositories.Content, repositories.Config, contentSearchClient)
	if err != nil {
		log.Fatalf("Failed to initialize content service: %v", err)
	}
	applicationServices[content.ServiceToken] = contentService

	contentSearchService, err := search.New(contentSearchClient)
	if err != nil {
		log.Fatalf("Failed to initialize search service: %v", err)
	}
	applicationServices[search.ContentSearchService] = contentSearchService

	uploadSearchService, err := search.New(uploadsSearchClient)
	if err != nil {
		log.Fatalf("Failed to initialize search service: %v", err)
	}
	applicationServices[search.UploadSearchService] = uploadSearchService

	storageService, err := storage.New(repositories.Uploads, repositories.Config, uploadsSearchClient, storageClient)
	if err != nil {
		log.Fatalf("Failed to initialize storage applicationServices: %v", err)
	}
	applicationServices[storage.ServiceToken] = storageService
	// End initialize applicationServices

	// Initialize Middlewares
	middlewares := make(middleware.Middlewares)
	CacheControlMiddleware := middleware.NewCacheControlMiddleware(repositories.Config)

	middlewares[middleware.CacheControlMiddleware] = CacheControlMiddleware
	middlewares[middleware.AnalyticsRecorderMiddleware] = middleware.NewAnalyticsRecorderMiddleware(analyticsService)
	middlewares[middleware.AuthMiddleware] = middleware.NewAuthMiddleware(conf.Paths, authService)
	middlewares[middleware.GzipMiddleware] = middleware.NewGzipMiddleware(configService)
	middlewares[middleware.CorsMiddleware] = middleware.NewCORSMiddleware(configService, CacheControlMiddleware)
	// End initialize middlewares

	// Initialize Handlers
	controllers.RegisterRoutes(conf.Paths, mux, applicationServices, middlewares)
	api.RegisterRoutes(mux, applicationServices, middlewares)
	// End Initialize Handlers

	return &application{
		mux:          mux,
		repositories: repositories,
		services:     applicationServices,
		middlewares:  middlewares,
		serveConfig:  conf.ServeConfig,
	}, nil
}
