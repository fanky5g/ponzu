package controllers

import (
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/domain/entities"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/views"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/services/config"
	"github.com/fanky5g/ponzu/internal/services/users"
	"github.com/fanky5g/ponzu/internal/util"
	"log"
	"net/http"
	"strings"
)

func hasSystemUsers(userService users.Service) (bool, error) {
	systemUsers, err := userService.ListUsers()
	if err != nil {
		return false, err
	}

	return len(systemUsers) > 0, nil
}

func NewLoginHandler(
	pathConf conf.Paths,
	configService config.Service,
	authService auth.Service,
	userService users.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		systemInitialized, err := hasSystemUsers(userService)
		if err != nil {
			log.Printf("Failed to check system initialization: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !systemInitialized {
			util.Redirect(req, res, pathConf, "/init", http.StatusFound)
			return
		}

		isValid, err := authService.IsTokenValid(request.GetAuthToken(req))
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Printf("Failed to check token validity: %v\n", err)
			return
		}

		if isValid {
			util.Redirect(req, res, pathConf, "/admin", http.StatusFound)
			return
		}

		switch req.Method {
		case http.MethodGet:
			appName, err := configService.GetAppName()
			if err != nil {
				log.Printf("Failed to get app name: %v\n", appName)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			view, err := views.Login(appName, pathConf)
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			res.Header().Set("Content-Type", "text/html")
			res.Write(view)

		case http.MethodPost:
			err = req.ParseForm()
			if err != nil {
				log.Println(err)
				util.Redirect(req, res, pathConf, req.URL.RequestURI(), http.StatusFound)
				return
			}

			email := strings.ToLower(req.FormValue("email"))
			password := req.FormValue("password")
			authToken, err := authService.LoginByEmail(email, &entities.Credential{
				Type:  entities.CredentialTypePassword,
				Value: password,
			})

			if err != nil || authToken == nil {
				log.Println("Failed to login user", err)
				util.Redirect(req, res, pathConf, req.URL.RequestURI(), http.StatusFound)
				return
			}

			http.SetCookie(res, &http.Cookie{
				Name:    "_token",
				Value:   authToken.Token,
				Expires: authToken.Expires,
				Path:    "/",
			})

			util.Redirect(req, res, pathConf, "/login", http.StatusFound)
		}
	}
}
