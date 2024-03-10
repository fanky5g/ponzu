package controllers

import (
	conf "github.com/fanky5g/ponzu/config"
	localStorage "github.com/fanky5g/ponzu/driver/storage"
	"github.com/fanky5g/ponzu/internal/handler/controllers/middleware"
	"github.com/fanky5g/ponzu/internal/services"
	"github.com/fanky5g/ponzu/internal/services/analytics"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/services/config"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/internal/services/search"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/internal/services/users"
	"log"
	"net/http"
)

func RegisterRoutes(
	pathConf conf.Paths,
	mux *http.ServeMux,
	services services.Services,
	middlewares middleware.Middlewares,
) {
	Auth := middlewares.Get(middleware.AuthMiddleware)
	CacheControlMiddleware := middleware.ToHttpHandler(middlewares.Get(middleware.CacheControlMiddleware))

	analyticsService := services.Get(analytics.ServiceToken).(analytics.Service)
	configService := services.Get(config.ServiceToken).(config.Service)
	userService := services.Get(users.ServiceToken).(users.Service)
	authService := services.Get(auth.ServiceToken).(auth.Service)
	contentService := services.Get(content.ServiceToken).(content.Service)
	storageService := services.Get(storage.ServiceToken).(storage.Service)
	contentSearchService := services.Get(search.ContentSearchService).(search.Service)
	uploadSearchService := services.Get(search.UploadSearchService).(search.Service)

	mux.HandleFunc("/", Auth(NewAdminHandler(pathConf, analyticsService, configService)))

	mux.HandleFunc("/init", NewInitHandler(pathConf, configService, userService, authService))

	mux.HandleFunc("/login", NewLoginHandler(pathConf, configService, authService, userService))
	mux.HandleFunc("/logout", NewLogoutHandler(pathConf))

	mux.HandleFunc("/recover", NewForgotPasswordHandler(pathConf, configService, userService, authService))
	mux.HandleFunc("/recover/key", NewRecoveryKeyHandler(pathConf, configService, authService, userService))

	mux.HandleFunc("/configure", Auth(NewConfigHandler(pathConf, configService)))
	mux.HandleFunc("/configure/users", Auth(NewConfigUsersHandler(pathConf, configService, authService, userService)))
	mux.HandleFunc("/configure/users/edit", Auth(NewConfigUsersEditHandler(pathConf, configService, authService, userService)))
	mux.HandleFunc("/configure/users/delete", Auth(NewConfigUsersDeleteHandler(pathConf, configService, authService, userService)))

	mux.HandleFunc("/uploads", Auth(NewUploadContentsHandler(pathConf, configService, storageService)))
	mux.HandleFunc("/uploads/search", Auth(NewUploadSearchHandler(pathConf, configService, uploadSearchService)))

	mux.HandleFunc("/contents", Auth(NewContentsHandler(pathConf, configService, contentService)))
	mux.HandleFunc("/contents/search", Auth(NewSearchHandler(pathConf, configService, contentSearchService)))
	mux.HandleFunc("/contents/export", Auth(NewExportHandler(pathConf, configService, contentService)))

	mux.HandleFunc("/edit", Auth(NewEditHandler(pathConf, configService, contentService, storageService)))
	mux.HandleFunc("/edit/delete", Auth(NewDeleteHandler(pathConf, configService, contentService)))
	mux.HandleFunc("/edit/upload", Auth(NewEditUploadHandler(pathConf, configService, storageService)))
	mux.HandleFunc("/edit/upload/delete", Auth(NewDeleteUploadHandler(pathConf, configService, storageService)))

	staticDir := conf.AdminStaticDir()
	staticFileSystem, err := localStorage.NewLocalStaticFileSystem(http.Dir(staticDir))
	if err != nil {
		log.Fatalf("Failed to create static file system: %v", err)
	}

	mux.Handle("/static/", CacheControlMiddleware(
		http.StripPrefix("/static", http.FileServer(staticFileSystem)),
	))
}
