package controllers

import (
	"net/http"

	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/api"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/references"
	"github.com/fanky5g/ponzu/internal/uploads"
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
	r.AuthorizedRoute("GET /contents/export", NewExportHandler)

	r.AuthorizedRoute("/edit/delete", NewDeleteHandler)
	r.AuthorizedRoute("/edit/upload/delete", NewDeleteUploadHandler)
	content.RegisterRoutes(r)
	uploads.RegisterRoutes(r)
	references.RegisterRoutes(r)

	api.RegisterRoutes(r, uploadsStaticFileSystem)

	r.HandleWithCacheControl("/static/", http.StripPrefix("/static", http.FileServer(staticFileSystem)))
	return nil
}
