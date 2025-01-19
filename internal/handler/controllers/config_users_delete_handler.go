package controllers

import (
	"net/http"
	"strings"

	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/services/users"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewConfigUsersDeleteHandler(r router.Router) http.HandlerFunc {
	authService := r.Context().Service(tokens.AuthServiceToken).(auth.Service)
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

			// do not allow current user to delete themselves
			user, err := authService.GetUserFromAuthToken(request.GetAuthToken(req))
			if err != nil {
				log.WithField("Error", err).Warning("Failed to get auth token")
				r.Renderer().InternalServerError(res)
				return
			}

			email := strings.ToLower(req.PostFormValue("email"))
			if user.Email == email {
				log.Debug("cannot delete own user account")
				r.Renderer().BadRequest(res)
				return
			}

			// delete existing user
			err = userService.DeleteUser(email)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to delete user")
				r.Renderer().InternalServerError(res)
				return
			}

			r.Redirect(req, res, strings.TrimSuffix(req.URL.RequestURI(), "/delete"))
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
