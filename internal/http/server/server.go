package server

import (
	"fmt"
	"github.com/fanky5g/ponzu/internal/storage"
	"net/http"

	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/internal/handler/controllers"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services"
	"github.com/fanky5g/ponzu/internal/services/analytics"
	"github.com/fanky5g/ponzu/internal/services/tls"
	"github.com/fanky5g/ponzu/tokens"
)

type Server interface {
	Serve() error
	ServeMux() *http.ServeMux
}

type server struct {
	tlsService       tls.Service
	configService    *config.Service
	analyticsService analytics.Service
	mux              *http.ServeMux
	configRepository database.Repository
}

func (server *server) ServeMux() *http.ServeMux {
	return server.mux
}

func New(
	contentTypes content.Types,
	assetStorage http.FileSystem,
	uploadStorage storage.Client,
	svcs services.Services,
	rootMux *http.ServeMux,
) (Server, error) {
	appConf, err := conf.Get()
	if err != nil {
		return nil, err
	}

	configService := svcs.Get(tokens.ConfigServiceToken).(*config.Service)
	analyticsService := svcs.Get(tokens.AnalyticsServiceToken).(analytics.Service)

	cfg, err := configService.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get application config: %v", err)
	}

	if cfg == nil {
		// initialize config
		cfg = &config.Config{}

		cfg.HTTPSPort = fmt.Sprintf("%d", appConf.ServeConfig.HttpsPort)
		cfg.HTTPPort = fmt.Sprintf("%d", appConf.ServeConfig.HttpPort)

		bind := appConf.ServeConfig.Bind
		// save the bound address the system is listening on so internal system can make
		// HTTP api calls while in dev or production w/o adding more cli flags
		if bind == "" {
			bind = "localhost"
		}
		cfg.BindAddress = bind

		if err = configService.SetConfig(cfg); err != nil {
			return nil, fmt.Errorf("failed to initialize config: %v", err)
		}
	}

	mux := http.NewServeMux()
	rtr, err := router.New(mux, appConf.Paths, svcs, contentTypes)
	if err != nil {
		return nil, err
	}

	rootMux.Handle(
		fmt.Sprintf("%s/", appConf.Paths.PublicPath),
		http.StripPrefix(appConf.Paths.PublicPath, mux),
	)

	err = controllers.RegisterRoutes(rtr, assetStorage, uploadStorage)
	if err != nil {
		return nil, err
	}

	return &server{
		tlsService:       svcs.Get(tokens.TLSServiceToken).(tls.Service),
		configService:    configService,
		analyticsService: analyticsService,
		mux:              rootMux,
	}, nil
}
