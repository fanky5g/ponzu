package controllers

import (
	"github.com/fanky5g/ponzu/internal/auth"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	authServicePkg "github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/services/users"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func NewLoginHandler(r router.Router) http.HandlerFunc {
	userService := r.Context().Service(tokens.UserServiceToken).(users.Service)
	authService := r.Context().Service(tokens.AuthServiceToken).(authServicePkg.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		systemUsers, err := userService.ListUsers()
		if err != nil {
			log.WithField("Error", err).Warning("Failed to list users")
			r.Renderer().InternalServerError(res)
			return
		}

		systemInitialized := len(systemUsers) > 0
		if !systemInitialized {
			r.Redirect(req, res, "/init")
			return
		}

		isValid, err := authService.IsTokenValid(request.GetAuthToken(req))
		if err != nil {
			log.WithField("Error", err).Printf("Failed to check token validity: %v\n", err)
			r.Renderer().InternalServerError(res)
			return
		}

		if isValid {
			r.Redirect(req, res, "/admin")
			return
		}

		switch req.Method {
		case http.MethodGet:
			r.Renderer().Render(res, "login_admin")

		case http.MethodPost:
			err = req.ParseForm()
			if err != nil {
				log.WithField("Error", err).Warning("Failed to parse form")
				r.Redirect(req, res, req.URL.RequestURI())
				return
			}

			email := strings.ToLower(req.FormValue("email"))
			password := req.FormValue("password")
			var authToken *auth.AuthToken
			authToken, err = authService.LoginByEmail(email, &auth.Credential{
				Type:  auth.CredentialTypePassword,
				Value: password,
			})

			if err != nil {
				log.WithField("Error", err).Warning("Failed to login user")
				r.Redirect(req, res, req.URL.RequestURI())
				return
			}

			if authToken == nil {
				r.Redirect(req, res, req.URL.RequestURI())
				return
			}

			http.SetCookie(res, &http.Cookie{
				Name:    "_token",
				Value:   authToken.Token,
				Expires: authToken.Expires,
				Path:    "/",
			})

			r.Redirect(req, res, "/login")
		}
	}
}
