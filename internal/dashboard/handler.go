package dashboard

import (
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/templates"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

// TODO(B.B): move these types to appropriate place
type Handler func(layoutTemplate *template.Template, viewModel *RootViewModel) http.HandlerFunc
type LayoutRouteHandler func(Handler) http.HandlerFunc

func GetTemplate() (*template.Template, error) {
	return template.New("dashboard").Funcs(templates.GlobFuncs).Parse(
		templates.Html("views/dashboard.gohtml", "views/app-frame.gohtml"),
	)
}

func NewHandler(r router.Router) (LayoutRouteHandler, error) {
	configCache := r.Context().Service(tokens.ConfigCache).(config.ConfigCache)

	layoutTemplate, templateErr := GetTemplate()
	if templateErr != nil {
		return nil, templateErr
	}

	return func(pageHandler Handler) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			rootViewModel, err := NewDashboardRootViewModel(
				configCache,
				r.Context().Paths().PublicPath,
				r.Context().Types().Content,
			)

			if err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Warning("Failed to create root view model")

				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			pageHandler(layoutTemplate, rootViewModel).ServeHTTP(res, req)
		}
	}, nil
}
