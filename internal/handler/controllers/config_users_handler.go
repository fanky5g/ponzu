package controllers

import (
	"net/http"
	"strings"

	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/services/users"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

func NewConfigUsersHandler(r router.Router) http.HandlerFunc {
	authService := r.Context().Service(tokens.AuthServiceToken).(auth.Service)
	userService := r.Context().Service(tokens.UserServiceToken).(users.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			currentUser, err := authService.GetUserFromAuthToken(request.GetAuthToken(req))
			if err != nil {
				log.WithField("Error", err).Warning("Failed to get user from request auth token")
				r.Renderer().InternalServerError(res)
				return
			}

			systemUsers, err := userService.ListUsers()
			if err != nil {
				log.WithField("Error", err).Warning("Failed to list users")
				r.Renderer().InternalServerError(res)
				return
			}

			for i, user := range systemUsers {
				if user.Email == currentUser.Email {
					systemUsers = append(systemUsers[:i], systemUsers[i+1:]...)
				}
			}

			r.Renderer().InjectTemplateInAdmin(res, r.Renderer().TemplateString("users_list"), map[string]interface{}{
				"GetUserByEmail": currentUser,
				"Users":          systemUsers,
			})

		case http.MethodPost:
			// create new user
			err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
			if err != nil {
				log.WithField("Error", err).Warning("Failed to parse form")
				r.Renderer().InternalServerError(res)
				return
			}

			email := strings.ToLower(req.FormValue("email"))
			password := req.PostFormValue("password")
			if email == "" || password == "" {
				log.Warn("received empty required values")
				r.Renderer().BadRequest(res)
				return
			}

			user, err := userService.CreateUser(email)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to create user")
				r.Renderer().InternalServerError(res)
				return
			}

			if err = authService.SetCredential(user.ID, &entities.Credential{
				Type:  entities.CredentialTypePassword,
				Value: password,
			}); err != nil {
				log.WithFields(
					log.Fields{
						"CredentialType": entities.CredentialTypePassword,
					},
				).Warning("Failed to update user credential")
				r.Renderer().InternalServerError(res)
				return
			}

			r.Redirect(req, res, req.URL.RequestURI())
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
