package controllers

import (
	"net/http"
	"strings"

	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	service "github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/auth"
	"github.com/fanky5g/ponzu/internal/services/users"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewConfigUsersEditHandler(r router.Router) http.HandlerFunc {
	authService := r.Context().Service(tokens.AuthServiceToken).(service.Service)
	userService := r.Context().Service(tokens.UserServiceToken).(users.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
			if err != nil {
				log.WithField("Error", err).Warning("Failed to parse form")
				r.Renderer().InternalServerError(res)
				return
			}

			user, err := authService.GetUserFromAuthToken(request.GetAuthToken(req))
			if err != nil {
				log.WithField("Error", err).Warning("Failed to get user from auth token")
				return
			}

			// check if password matches
			password := &auth.Credential{
				Type:  auth.CredentialTypePassword,
				Value: req.PostFormValue("password"),
			}

			if err = authService.VerifyCredential(user.ID, password); err != nil {
				log.WithField("Error", err).Warning("Unexpected user/password combination")
				r.Renderer().BadRequest(res)
				return
			}

			email := strings.ToLower(req.PostFormValue("email"))
			newPassword := req.PostFormValue("new_password")
			if newPassword != "" {
				if err = authService.SetCredential(user.ID, &auth.Credential{
					Type:  auth.CredentialTypePassword,
					Value: newPassword,
				}); err != nil {
					log.WithField("Error", err).Warning("Failed to update password")
					r.Renderer().InternalServerError(res)
					return
				}
			}

			if email != "" {
				update := &auth.User{
					ID:    user.ID,
					Email: email,
				}

				if err = userService.UpdateUser(user, update); err != nil {
					log.WithField("Error", err).Warning("Failed to update user")
					r.Renderer().InternalServerError(res)
					return
				}

				user = update
			}

			// create new token
			authToken, err := authService.NewToken(user)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to generate token")
				r.Renderer().InternalServerError(res)
				return
			}

			cookie := &http.Cookie{
				Name:    "_token",
				Value:   authToken.Token,
				Expires: authToken.Expires,
				Path:    "/",
			}

			http.SetCookie(res, cookie)
			// add new token cookie to the request
			req.AddCookie(cookie)
			r.Redirect(req, res, strings.TrimSuffix(req.URL.RequestURI(), "/edit"))

		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
