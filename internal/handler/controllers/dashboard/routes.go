package dashboard

import (
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/internal/handler/controllers/dashboard/content/edit"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
)

func RegisterRoutes(
	r router.Router,
	staticFileSystem driver.StaticFileSystemInterface,
	uploadsStaticFileSystem driver.StaticFileSystemInterface,
) error {
	r.AuthorizedRoute("/edit", edit.Handler)

	return nil
}
