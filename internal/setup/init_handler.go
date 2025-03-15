package setup

import (
	"encoding/base64"
	"fmt"
	"github.com/fanky5g/ponzu/internal/auth"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type ViewModel struct {
	PublicPath string
}

func NewInitPageHandler(publicPath string, userService *auth.UserService, layout layouts.Template) http.HandlerFunc {
	tmpl, templateErr := layout.Child("views/init_admin.gohtml")
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
		if systemInitialized {
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
			response.NewHTMLResponse(http.StatusOK, tmpl, ViewModel{PublicPath: publicPath}),
		)
	}
}

func NewInitHandler(publicPath string, configService *config.Service, authService *auth.Service, userService *auth.UserService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			log.WithField("Error", err).Warning("Failed to parse form")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		cfg, err := configService.Get()
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get config")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if cfg.ClientSecret == "" {
			clientSecret := fmt.Sprintf("%s%d", req.FormValue("name"), time.Now().Unix())
			cfg.ClientSecret = base64.StdEncoding.EncodeToString([]byte(clientSecret))
		}

		// create and save controllers user
		email := strings.ToLower(req.FormValue("email"))
		password := req.FormValue("password")
		appName := req.FormValue("name")

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

		cfg.Name = appName

		var authToken *auth.Token
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

		response.Respond(
			res,
			req,
			response.NewRedirectResponse(publicPath, "/init"),
		)
	}
}
