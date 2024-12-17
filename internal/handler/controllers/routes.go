package controllers

import (
	"net/http"

	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/internal/handler/controllers/api"
	"github.com/fanky5g/ponzu/internal/handler/controllers/dashboard"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
)

func RegisterRoutes(
	r router.Router,
	staticFileSystem driver.StaticFileSystemInterface,
	uploadsStaticFileSystem driver.StaticFileSystemInterface,
) error {
	r.AuthorizedRoute("/", NewAdminHandler)

	r.Route("/init", NewInitHandler)
	r.Route("/login", NewLoginHandler)
	r.Route("/logout", NewLogoutHandler)
	r.Route("/recover", NewForgotPasswordHandler)
	r.Route("/recover/key", NewRecoveryKeyHandler)

	r.AuthorizedRoute("/configure", NewConfigHandler)
	r.AuthorizedRoute("/configure/users", NewConfigUsersHandler)
	r.AuthorizedRoute("/configure/users/edit", NewConfigUsersEditHandler)
	r.AuthorizedRoute("/configure/users/delete", NewConfigUsersDeleteHandler)

	r.AuthorizedRoute("/uploads", NewUploadContentsHandler)
	r.AuthorizedRoute("/uploads/search", NewUploadSearchHandler)
	r.AuthorizedRoute("/contents", NewContentsHandler)
	r.AuthorizedRoute("/contents/search", NewSearchHandler)
	r.AuthorizedRoute("/contents/export", NewExportHandler)

	r.AuthorizedRoute("/edit/delete", NewDeleteHandler)
	r.AuthorizedRoute("/edit/upload", NewEditUploadHandler)
	r.AuthorizedRoute("/edit/upload/delete", NewDeleteUploadHandler)
	dashboard.RegisterRoutes(r, staticFileSystem, uploadsStaticFileSystem)

	api.RegisterRoutes(r, uploadsStaticFileSystem)

	r.HandleWithCacheControl("/static/", http.StripPrefix("/static", http.FileServer(staticFileSystem)))
	return nil
}
