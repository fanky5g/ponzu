package controllers

import (
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func NewForgotPasswordHandler(r router.Router) http.HandlerFunc {
	authService := r.Context().Service(tokens.AuthServiceToken).(auth.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			r.Renderer().Render(res, "forgot_password")

		case http.MethodPost:
			err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
			if err != nil {
				log.WithField("Error", err).Warning("Failed to parse form")
				r.Renderer().InternalServerError(res)
				return
			}

			// check email for user, if no user return Error
			email := strings.ToLower(req.FormValue("email"))
			if email == "" {
				log.Info("Failed account recovery. No email address submitted.")
				r.Renderer().BadRequest(res)
				return
			}

			if err = authService.SendPasswordRecoveryInstructions(email); err != nil {
				log.WithField("Error", err).Warning("Failed to send password recovery instructions")
			}

			r.Redirect(req, res, "/recover/key")
		default:
			r.Renderer().MethodNotAllowed(res)
			return
		}
	}
}
