package server

import (
	"fmt"
	"github.com/fanky5g/ponzu/internal/analytics"
	"github.com/fanky5g/ponzu/internal/auth"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/content/dataexporter"
	"github.com/fanky5g/ponzu/internal/content/dataexporter/formatter"
	"github.com/fanky5g/ponzu/internal/http/middleware"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/http/router"
	"github.com/fanky5g/ponzu/internal/layouts"
	"github.com/fanky5g/ponzu/internal/layouts/dashboard"
	"github.com/fanky5g/ponzu/internal/layouts/root"
	"github.com/fanky5g/ponzu/internal/memorycache"
	"github.com/fanky5g/ponzu/internal/search"
	"github.com/fanky5g/ponzu/internal/setup"

	"github.com/fanky5g/ponzu/internal/storage"
	"github.com/pkg/errors"
	"net/http"

	conf "github.com/fanky5g/ponzu/config"
	contentPkg "github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/database"
)

type Server interface {
	ServeMux() *http.ServeMux
}

type server struct {
	configService    *config.Service
	mux              *http.ServeMux
	configRepository database.Repository
}

func (server *server) ServeMux() *http.ServeMux {
	return server.mux
}

func New(
	contentTypes map[string]contentPkg.Builder,
	db database.Database,
	assetStorage http.FileSystem,
	uploadStorage storage.Client,
	searchClient search.SearchInterface,
	rootMux *http.ServeMux,
) (Server, error) {
	appConf, err := conf.Get()
	if err != nil {
		return nil, err
	}

	memcache, err := memorycache.New()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize memory cache")
	}

	configCache, err := config.NewCache(memcache)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create config cache service")
	}

	configService, err := config.New(db, memcache)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create config service")
	}

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

	rowFormatter, err := formatter.NewJSONRowFormatter()
	if err != nil {
		return nil, err
	}

	contentExporter, err := dataexporter.New(rowFormatter)
	if err != nil {
		return nil, err
	}

	storageService, err := content.NewUploadService(db, searchClient, uploadStorage)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize storage service")
	}

	contentService, err := content.New(db, contentTypes, searchClient, contentExporter, storageService)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize content service")
	}

	uploadService, err := content.NewUploadService(db, searchClient, uploadStorage)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize upload service")
	}

	searchService, err := search.New(searchClient, db)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize search service")
	}

	clientSecret, err := configService.GetStringValue("client_secret")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get client secret")
	}

	authService, err := auth.New(clientSecret, db)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create auth service")
	}

	userService, err := auth.NewUserService(db)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create user service")
	}

	analyticsService, err := analytics.New(db)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize analytics services")
	}

	middlewares, err := middleware.New(appConf.Paths, authService, configCache, analyticsService)

	mux := http.NewServeMux()
	rtr, err := router.New(mux, middlewares)
	if err != nil {
		return nil, err
	}

	dashboardLayout, err := dashboard.New(configCache, appConf.Paths.PublicPath, contentTypes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create dashboard layout")
	}

	rootLayout, err := root.New(configCache, appConf.Paths.PublicPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create dashboard layout")
	}

	setup.RegisterRoutes(rtr, configService, authService, userService, appConf.Paths.PublicPath, rootLayout)
	analytics.RegisterRoutes(rtr, analyticsService, dashboardLayout)
	auth.RegisterRoutes(rtr, authService, userService, appConf.Paths.PublicPath, rootLayout, dashboardLayout)
	config.RegisterRoutes(rtr, appConf.Paths.PublicPath, configService, dashboardLayout)

	if err = content.RegisterRoutes(
		rtr,
		appConf.Paths.PublicPath,
		contentTypes,
		contentService,
		uploadService,
		searchService,
		dashboardLayout,
	); err != nil {
		return nil, err
	}

	rtr.HandleWithCacheControl("GET /static/", http.StripPrefix("/static", http.FileServer(assetStorage)))
	rtr.HandleWithCacheControl(
		"GET /api/uploads/",
		http.StripPrefix("/api/uploads/", http.FileServer(uploadStorage)),
	)

	// Catch-All Route - 404 page
	var error404 layouts.Template
	error404, err = rootLayout.Child("views/errors/error_404.gohtml")
	if err != nil {
		return nil, err
	}

	rtr.Route("GET /", func() http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			response.Respond(
				res,
				req,
				response.NewHTMLResponse(
					http.StatusNotFound,
					error404,
					nil,
				),
			)
		}
	})
	// End Catch-All Route

	rootMux.Handle(
		fmt.Sprintf("%s/", appConf.Paths.PublicPath),
		http.StripPrefix(appConf.Paths.PublicPath, mux),
	)

	return &server{
		configService: configService,
		mux:           rootMux,
	}, nil
}
