package controllers

import (
	"encoding/base64"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	authServicePkg "github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/auth"
	"github.com/fanky5g/ponzu/internal/services/users"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/fanky5g/ponzu/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func NewInitHandler(r router.Router) http.HandlerFunc {
	configService := r.Context().Service(tokens.ConfigServiceToken).(*config.Service)
	authService := r.Context().Service(tokens.AuthServiceToken).(authServicePkg.Service)
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

			var cfg *config.Config
			cfg, err = configService.Get()
			if err != nil {
				log.WithField("Error", err).Warning("Failed to get config")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			if cfg.ClientSecret == "" {
				name := []byte(req.FormValue("name") + util.NewEtag())
				cfg.ClientSecret = base64.StdEncoding.EncodeToString(name)
			}

			// generate an Etag to use for response caching
			cfg.Etag = util.NewEtag()

			// create and save controllers user
			email := strings.ToLower(req.FormValue("email"))
			password := req.FormValue("password")
			var user *auth.User
			user, err = userService.CreateUser(email)
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err = authService.SetCredential(user.ID, &auth.Credential{
				Type:  auth.CredentialTypePassword,
				Value: password,
			}); err != nil {
				log.WithField("Error", err).Warning("Failed to create admin user")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			// set initial user email as admin_email and make config
			cfg.AdminEmail = email
			err = configService.SetConfig(cfg)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to update config")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			var authToken *auth.AuthToken
			authToken, err = authService.NewToken(user)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to generate auth token")
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
