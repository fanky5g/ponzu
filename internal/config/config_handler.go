package config

import (
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type ViewModel struct {
	PublicPath string
	Form       template.HTML
}

func NewEditConfigHandler(publicPath string, configService *Service, layout layouts.Template) http.HandlerFunc {
	tmpl, templateErr := layout.Child("views/edit_config_view.gohtml")
	if templateErr != nil {
		panic(templateErr)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		cfg, err := configService.Get()
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		form, err := cfg.MarshalEditor(publicPath)
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewHTMLResponse(
				http.StatusOK,
				tmpl,
				ViewModel{
					PublicPath: publicPath,
					Form:       template.HTML(form),
				},
			),
		)
	}
}

func NewSaveConfigHandler(configService *Service, publicPath string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		entity, err := request.GetEntity(Builder(), req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map config entity")
			return
		}

		err = configService.SetConfig(entity.(*Config))
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(res, req, response.NewRedirectResponse(publicPath, req.URL.RequestURI()))
	}
}
