package auth

import (
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func NewLoginHandler(publicPath string, authService *Service, userService *UserService, layout layouts.Template) http.HandlerFunc {
	tmpl, templateErr := layout.Child("views/login_admin")
	if templateErr != nil {
		panic(templateErr)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		systemUsers, err := userService.ListUsers()
		if err != nil {
			log.WithField("Error", err).Warning("Failed to list users")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		systemInitialized := len(systemUsers) > 0
		if !systemInitialized {
			response.Respond(
				res,
				req,
				response.NewRedirectResponse(publicPath, "/init"),
			)
			return
		}

		isValid, err := authService.IsTokenValid(request.GetAuthToken(req))
		if err != nil {
			log.WithField("Error", err).Printf("Failed to check token validity: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if isValid {
			response.Respond(
				res,
				req,
				response.NewRedirectResponse(publicPath, "/"),
			)
			return
		}

		response.Respond(
			res,
			req,
			response.NewHTMLResponse(http.StatusOK, tmpl, &ViewModel{PublicPath: publicPath}),
		)
	}
}

func NewAuthHandler(publicPath string, authService *Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse form")
			response.Respond(
				res,
				req,
				response.NewRedirectResponse(publicPath, req.URL.RequestURI()),
			)
			return
		}

		email := strings.ToLower(req.FormValue("email"))
		password := req.FormValue("password")
		var authToken *Token
		authToken, err = authService.LoginByEmail(email, &Credential{
			Type:  CredentialTypePassword,
			Value: password,
		})

		if err != nil {
			log.WithField("Error", err).Warning("Failed to login user")
			response.Respond(
				res,
				req,
				response.NewRedirectResponse(publicPath, req.URL.RequestURI()),
			)
			return
		}

		if authToken == nil {
			response.Respond(
				res,
				req,
				response.NewRedirectResponse(publicPath, req.URL.RequestURI()),
			)
			return
		}

		http.SetCookie(res, &http.Cookie{
			Name:    "_token",
			Value:   authToken.Token,
			Expires: authToken.Expires,
			Path:    "/",
		})

		response.Respond(
			res,
			req,
			response.NewRedirectResponse(publicPath, "/"),
		)
	}
}
