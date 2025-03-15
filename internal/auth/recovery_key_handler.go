package auth

import (
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func NewRecoveryKeyPageHandler(publicPath string, layout layouts.Template) http.HandlerFunc {
	tmpl, err := layout.Child("views/recovery_key.gohtml")
	if err != nil {
		panic(err)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		response.Respond(
			res,
			req,
			response.NewHTMLResponse(http.StatusOK, tmpl, &ViewModel{PublicPath: publicPath}),
		)
	}
}

func NewRecoveryKeyHandler(publicPath string, authService *Service, userService *UserService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
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
		var recoveryKey *RecoveryKey
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
		var user *User
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

		if err = authService.SetCredential(user.ID, &Credential{
			Type:  CredentialTypePassword,
			Value: password,
		}); err != nil {
			log.WithField("Error", err).Warning("Error updating user")

			res.WriteHeader(http.StatusInternalServerError)
			if _, err = res.Write([]byte("Error, please go back and try again.")); err != nil {
				log.WithField("Error", err).Warning("Failed to write response")
			}
			return
		}

		response.Respond(
			res,
			req,
			response.NewRedirectResponse(publicPath, "/login"),
		)
	}
}
