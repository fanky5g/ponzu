package controllers

import (
	conf "github.com/fanky5g/ponzu/config"
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

func NewConfigUsersDeleteHandler(
	pathConf conf.Paths,
	configService config.Service,
	authService auth.Service,
	userService users.Service,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		appName, err := configService.GetAppName()
		if err != nil {
			log.Printf("Failed to get app name: %v\n", appName)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch req.Method {
		case http.MethodPost:
			err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
			if err != nil {
				LogAndFail(res, err, appName, pathConf)
				return
			}

			// do not allow current user to delete themselves
			user, err := authService.GetUserFromAuthToken(request.GetAuthToken(req))
			if err != nil {
				LogAndFail(res, err, appName, pathConf)
				return
			}

			email := strings.ToLower(req.PostFormValue("email"))
			if user.Email == email {
				log.Println(err)
				res.WriteHeader(http.StatusBadRequest)
				errView, err := views.Admin(util.Html("error_405"), appName, pathConf)
				if err != nil {
					return
				}

				res.Write(errView)
				return
			}

			// delete existing user
			err = userService.DeleteUser(email)
			if err != nil {
				LogAndFail(res, err, appName, pathConf)
				return
			}

			util.Redirect(req, res, pathConf, strings.TrimSuffix(req.URL.RequestURI(), "/delete"), http.StatusFound)

		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
