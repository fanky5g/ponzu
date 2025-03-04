package controllers

import (
	"github.com/fanky5g/ponzu/internal/auth"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	authServicePkg "github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/services/users"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func NewRecoveryKeyHandler(r router.Router) http.HandlerFunc {
	authService := r.Context().Service(tokens.AuthServiceToken).(authServicePkg.Service)
	userService := r.Context().Service(tokens.UserServiceToken).(users.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			r.Renderer().Render(res, "recovery_key")

		case http.MethodPost:
			err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
			if err != nil {
				log.WithField("Error", err).Warning("Failed to parse form")
				res.WriteHeader(http.StatusInternalServerError)
				if _, err = res.Write([]byte("Error, please go back and try again.")); err != nil {
					log.WithField("Error", err).Warning("Failed to write response")
				}

				return
			}

			// check for email & key match
			email := strings.ToLower(req.FormValue("email"))
			key := req.FormValue("key")

			var actual string
			var recoveryKey *auth.RecoveryKey
			recoveryKey, err = authService.GetRecoveryKey(email)
			if err != nil {
				log.WithField("Error", err).Warning("Error getting recovery key from database")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			if recoveryKey == nil {
				if _, err = res.Write([]byte("Error, please go back and try again.")); err != nil {
					log.WithField("Error", err).Warning("Failed to write response")
				}

				res.WriteHeader(http.StatusBadRequest)
				return
			}

			actual = recoveryKey.Value
			if key != actual {
				log.WithField("key", key).Warning("Bad recovery key submitted")
				res.WriteHeader(http.StatusBadRequest)
				if _, err = res.Write([]byte("Error, please go back and try again.")); err != nil {
					log.WithField("Error", err).Warning("Failed to write response")
				}

				return
			}

			// set user with new password
			password := req.FormValue("password")
			var user *auth.User
			user, err = userService.GetUserByEmail(email)
			if err != nil {
				log.WithField("Error", err).Warning("Error finding user by email")
				res.WriteHeader(http.StatusBadRequest)
				if _, err = res.Write([]byte("Error, please go back and try again.")); err != nil {
					log.WithField("Error", err).Warning("Failed to write response")
				}

				return
			}

			if user == nil {
				log.Warning("No user found with email")
				res.WriteHeader(http.StatusBadRequest)
				if _, err = res.Write([]byte("Error, please go back and try again.")); err != nil {
					log.WithField("Error", err).Warning("Failed to write response")
				}

				return
			}

			if err = authService.SetCredential(user.ID, &auth.Credential{
				Type:  auth.CredentialTypePassword,
				Value: password,
			}); err != nil {
				log.WithField("Error", err).Warning("Error updating user")

				res.WriteHeader(http.StatusInternalServerError)
				if _, err = res.Write([]byte("Error, please go back and try again.")); err != nil {
					log.WithField("Error", err).Warning("Failed to write response")
				}
				return
			}

			r.Redirect(req, res, "/login")
			return
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
