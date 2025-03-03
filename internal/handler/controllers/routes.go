package controllers

import (
	"github.com/fanky5g/ponzu/internal/dashboard"
	"net/http"

	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/api"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/storage"
)

func RegisterRoutes(r router.Router, assetStorage http.FileSystem, uploadStorage storage.Client) error {
	dashboardHandler, err := dashboard.NewHandler(r)
	if err != nil {
		return err
	}

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
	err = content.RegisterRoutes(r, dashboardHandler)
	if err != nil {
		return err
	}

	api.RegisterRoutes(r, uploadStorage)

	r.HandleWithCacheControl("/static/", http.StripPrefix("/static", http.FileServer(assetStorage)))
	return nil
}
