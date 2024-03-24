package controllers

import (
	"encoding/json"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/config"
	"github.com/fanky5g/ponzu/tokens"
	"log"
	"net/http"
)

func NewConfigHandler(r router.Router) http.HandlerFunc {
	configService := r.Context().Service(tokens.ConfigServiceToken).(config.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			data, err := configService.GetAll()
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			c := &entities.Config{}

			err = json.Unmarshal(data, c)
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Renderer().Editable(res, c)
		case http.MethodPost:
			err := req.ParseForm()
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = configService.SetConfig(req.Form)
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
