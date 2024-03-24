package controllers

import (
	"encoding/base64"
	entities2 "github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/services/config"
	"github.com/fanky5g/ponzu/internal/services/users"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/fanky5g/ponzu/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func NewInitHandler(r router.Router) http.HandlerFunc {
	configService := r.Context().Service(tokens.ConfigServiceToken).(config.Service)
	authService := r.Context().Service(tokens.AuthServiceToken).(auth.Service)
	userService := r.Context().Service(tokens.UserServiceToken).(users.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		systemUsers, err := userService.ListUsers()
		if err != nil {
			log.WithField("Error", err).Warning("Failed to list users")
			r.Renderer().InternalServerError(res)
			return
		}

		systemInitialized := len(systemUsers) > 0
		if systemInitialized {
			r.Redirect(req, res, "/admin")
			return
		}

		switch req.Method {
		case http.MethodGet:
			r.Renderer().Render(res, "init_admin")
			return
		case http.MethodPost:
			err = req.ParseForm()
			if err != nil {
				log.WithField("Error", err).Warning("Failed to parse form")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			// get the site name from post to encode and use as secret
			name := []byte(req.FormValue("name") + util.NewEtag())
			secret := base64.StdEncoding.EncodeToString(name)
			req.Form.Set("client_secret", secret)

			// generate an Etag to use for response caching
			etag := util.NewEtag()
			req.Form.Set("etag", etag)

			// create and save controllers user
			email := strings.ToLower(req.FormValue("email"))
			password := req.FormValue("password")
			var user *entities2.User
			user, err = userService.CreateUser(email)
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err = authService.SetCredential(user.ID, &entities2.Credential{
				Type:  entities2.CredentialTypePassword,
				Value: password,
			}); err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			// set HTTP port which should be previously added to config cache
			var port string
			port, err = configService.GetCacheStringValue("http_port")
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			req.Form.Set("http_port", port)
			// set initial user email as admin_email and make config
			req.Form.Set("admin_email", email)
			err = configService.SetConfig(req.Form)
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			authToken, err := authService.NewToken(user)
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(res, &http.Cookie{
				Name:    "_token",
				Value:   authToken.Token,
				Expires: authToken.Expires,
				Path:    "/",
			})

			r.Redirect(req, res, strings.TrimSuffix(req.URL.RequestURI(), "/init"))
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

}
