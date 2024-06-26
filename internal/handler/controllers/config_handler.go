package controllers

import (
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/config"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewConfigHandler(r router.Router) http.HandlerFunc {
	configService := r.Context().Service(tokens.ConfigServiceToken).(config.Service)

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

			entity, err := request.GetEntityFromFormData(entities.ConfigBuilder, req.PostForm)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to map config entity")
				return
			}

			err = configService.SetConfig(entity.(*entities.Config))
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
