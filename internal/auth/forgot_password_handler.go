package auth

import (
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type ViewModel struct {
	PublicPath string
}

func NewForgotPasswordPageHandler(publicPath string, layout layouts.Template) http.HandlerFunc {
	tmpl, err := layout.Child("views/forgot_password.gohtml")
	if err != nil {
		panic(err)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		response.Respond(
			res,
			req,
			response.NewHTMLResponse(http.StatusOK, tmpl, ViewModel{PublicPath: publicPath}),
		)
	}
}

func NewForgotPasswordHandler(publicPath string, authService *Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse form")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// check email for user, if no user return Error
		email := strings.ToLower(req.FormValue("email"))
		if email == "" {
			log.Info("Failed account recovery. No email address submitted.")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = authService.SendPasswordRecoveryInstructions(email); err != nil {
			log.WithField("Error", err).Warning("Failed to send password recovery instructions")
		}

		response.Respond(
			res,
			req,
			response.NewRedirectResponse(publicPath, "/recover/key"),
		)
	}
}
