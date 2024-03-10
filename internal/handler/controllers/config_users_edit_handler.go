package controllers

import (
	"fmt"
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

func NewConfigUsersEditHandler(
	pathConf conf.Paths,
	configService config.Service,
	authService auth.Service,
	userService users.Service,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			appName, err := configService.GetAppName()
			if err != nil {
				log.Printf("Failed to get app name: %v\n", appName)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				errView, err := views.Admin(util.Html("error_500"), appName, pathConf)
				if err != nil {
					return
				}

				res.Write(errView)
				return
			}

			user, err := authService.GetUserFromAuthToken(request.GetAuthToken(req))
			if err != nil {
				LogAndFail(res, err, appName, pathConf)
				return
			}

			// check if password matches
			password := &entities.Credential{
				Type:  entities.CredentialTypePassword,
				Value: req.PostFormValue("password"),
			}

			if err = authService.VerifyCredential(user.ID, password); err != nil {
				log.Printf("Unexpected user/password combination: %v\n", err)
				res.WriteHeader(http.StatusBadRequest)
				errView, err := views.Admin(util.Html("error_405"), appName, pathConf)
				if err != nil {
					return
				}

				res.Write(errView)
				return
			}

			email := strings.ToLower(req.PostFormValue("email"))
			newPassword := req.PostFormValue("new_password")
			if newPassword != "" {
				if err = authService.SetCredential(user.ID, &entities.Credential{
					Type:  entities.CredentialTypePassword,
					Value: newPassword,
				}); err != nil {
					LogAndFail(res, fmt.Errorf("failed to update password: %v", err), appName, pathConf)
					return
				}
			}

			if email != "" {
				update := &entities.User{
					ID:    user.ID,
					Email: email,
				}

				if err = userService.UpdateUser(user, update); err != nil {
					LogAndFail(res, fmt.Errorf("failed to update user: %v", err), appName, pathConf)
					return
				}

				user = update
			}

			// create new token
			authToken, err := authService.NewToken(user)
			if err != nil {
				LogAndFail(res, fmt.Errorf("failed to generate token: %v", err), appName, pathConf)
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
			util.Redirect(req, res, pathConf, strings.TrimSuffix(req.URL.RequestURI(), "/edit"), http.StatusFound)

		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
