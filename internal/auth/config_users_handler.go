package auth

import (
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type UsersViewModel struct {
	PublicPath     string
	GetUserByEmail *User
	Users          []User
}

func NewUsersPageHandler(publicPath string, authService *Service, userService *UserService, layout layouts.Template) http.HandlerFunc {
	tmpl, templateErr := layout.Child("views/users_list.gohtml")
	if templateErr != nil {
		panic(templateErr)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		currentUser, err := authService.GetUserFromAuthToken(request.GetAuthToken(req))
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get user from request auth token")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		systemUsers, err := userService.ListUsers()
		if err != nil {
			log.WithField("Error", err).Warning("Failed to list users")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i, user := range systemUsers {
			if user.Email == currentUser.Email {
				systemUsers = append(systemUsers[:i], systemUsers[i+1:]...)
			}
		}

		response.Respond(
			res,
			req,
			response.NewHTMLResponse(http.StatusOK, tmpl, UsersViewModel{
				PublicPath:     publicPath,
				GetUserByEmail: currentUser,
				Users:          systemUsers,
			}),
		)
	}
}

func NewCreateUserHandler(publicPath string, authService *Service, userService *UserService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// create new user
		err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse form")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		email := strings.ToLower(req.FormValue("email"))
		password := req.PostFormValue("password")
		if email == "" || password == "" {
			log.Warn("received empty required values")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := userService.CreateUser(email)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to create user")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = authService.SetCredential(user.ID, &Credential{
			Type:  CredentialTypePassword,
			Value: password,
		}); err != nil {
			log.WithFields(
				log.Fields{
					"CredentialType": CredentialTypePassword,
				},
			).Warning("Failed to update user credential")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewRedirectResponse(publicPath, req.URL.RequestURI()),
		)
	}
}
