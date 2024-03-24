package server

import (
	"fmt"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/infrastructure"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
	"github.com/fanky5g/ponzu/internal/handler/controllers"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services"
	"github.com/fanky5g/ponzu/internal/services/tls"
	"github.com/fanky5g/ponzu/tokens"
	"net/http"
)

type Server interface {
	Serve() error
	ServeMux() *http.ServeMux
}

type server struct {
	cfg         *config.Config
	tlsService  tls.Service
	configCache repositories.Cache
	mux         *http.ServeMux
}

func (server *server) ServeMux() *http.ServeMux {
	return server.mux
}

func New(contentTypes content.Types, infra infrastructure.Infrastructure, svcs services.Services) (Server, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	db := infra.Get(tokens.DatabaseInfrastructureToken).(driver.Database)

	configRepository := db.Get(tokens.ConfigRepositoryToken).(repositories.ConfigRepositoryInterface)
	err = configRepository.PutConfig("https_port", fmt.Sprintf("%d", cfg.ServeConfig.HttpsPort))
	if err != nil {
		return nil, fmt.Errorf("failed to save server config: %v", err)
	}

	// save the https port the system is listening on so internal system can make
	// HTTP api calls while in dev or production w/o adding more cli flags
	err = configRepository.PutConfig("http_port", fmt.Sprintf("%d", cfg.ServeConfig.HttpPort))
	if err != nil {
		return nil, fmt.Errorf("failed to save server config: %v", err)
	}

	bind := cfg.ServeConfig.Bind
	// save the bound address the system is listening on so internal system can make
	// HTTP api calls while in dev or production w/o adding more cli flags
	if bind == "" {
		bind = "localhost"
	}

	err = configRepository.PutConfig("bind_addr", bind)
	if err != nil {
		return nil, fmt.Errorf("failed to save server config: %v", err)
	}

	mux := http.NewServeMux()
	configCache := configRepository.Cache()

	rtr, err := router.New(mux, cfg.Paths, configCache, svcs, contentTypes)
	if err != nil {
		return nil, err
	}

	err = controllers.RegisterRoutes(
		rtr,
		infra.Get(tokens.AssetStorageClientInfrastructureToken).(driver.StorageClientInterface),
		infra.Get(tokens.StorageClientInfrastructureToken).(driver.StorageClientInterface),
	)
	if err != nil {
		return nil, err
	}

	return &server{
		cfg:         cfg,
		tlsService:  svcs.Get(tokens.TLSServiceToken).(tls.Service),
		configCache: configCache,
		mux:         mux,
	}, nil
}
