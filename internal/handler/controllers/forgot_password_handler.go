package controllers

import (
	"errors"
	"fmt"
	domainErrors "github.com/fanky5g/ponzu/errors"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/services/config"
	"github.com/fanky5g/ponzu/internal/services/users"
	"github.com/fanky5g/ponzu/tokens"
	emailer "github.com/nilslice/email"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func NewForgotPasswordHandler(r router.Router) http.HandlerFunc {
	userService := r.Context().Service(tokens.UserServiceToken).(users.Service)
	authService := r.Context().Service(tokens.AuthServiceToken).(auth.Service)
	configService := r.Context().Service(tokens.ConfigServiceToken).(config.Service)

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

			_, err = userService.GetUserByEmail(email)
			if errors.Is(err, domainErrors.ErrNoUserExists) {
				log.WithField("Error", err).Info("No user exists")
				r.Renderer().BadRequest(res)
				return
			}

			// create temporary key to verify user
			key, err := authService.SetRecoveryKey(email)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to set account recovery key.", err)
				return
			}

			domain, err := configService.GetStringValue("domain")
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to get domain from configuration.", err)
				return
			}

			body := fmt.Sprintf(`
There has been an account recovery request made for the user with email:
%s

To recover your account, please go to http://%s/recover/key and enter 
this email address along with the following secret key:

%s

If you did not make the request, ignore this message and your password 
will remain as-is.


Thank you,
Ponzu CMS at %s

`, email, domain, key, domain)

			msg := emailer.Message{
				To:      email,
				From:    fmt.Sprintf("ponzu@%s", domain),
				Subject: fmt.Sprintf("Account Recovery [%s]", domain),
				Body:    body,
			}

			go func() {
				err = msg.Send()
				if err != nil {
					log.Println("Failed to send message to:", msg.To, "about", msg.Subject, "Error:", err)
				}
			}()

			r.Redirect(req, res, "/recover/key")
		default:
			r.Renderer().MethodNotAllowed(res)
			return
		}
	}
}
