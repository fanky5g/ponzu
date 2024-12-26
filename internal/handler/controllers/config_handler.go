package controllers

import (
	"net/http"

	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewConfigHandler(r router.Router) http.HandlerFunc {
	configService := r.Context().Service(tokens.ConfigServiceToken).(*config.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			cfg, err := configService.Get()
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Renderer().Editable(res, cfg)
		case http.MethodPost:
			err := req.ParseForm()
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			entity, err := request.GetEntityFromFormData(config.ConfigBuilder, req.PostForm)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to map config entity")
				return
			}

			err = configService.SetConfig(entity.(*config.Config))
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Redirect(req, res, req.URL.RequestURI())
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
