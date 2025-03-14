package auth

import (
	"github.com/fanky5g/ponzu/internal/http/response"
	"net/http"
	"strings"

	"github.com/fanky5g/ponzu/internal/http/request"
	log "github.com/sirupsen/logrus"
)

func NewConfigUsersEditHandler(publicPath string, authService *Service, userService *UserService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse form")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := authService.GetUserFromAuthToken(request.GetAuthToken(req))
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get user from auth token")
			return
		}

		// check if password matches
		password := &Credential{
			Type:  CredentialTypePassword,
			Value: req.PostFormValue("password"),
		}

		if err = authService.VerifyCredential(user.ID, password); err != nil {
			log.WithField("Error", err).Warning("Unexpected user/password combination")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		email := strings.ToLower(req.PostFormValue("email"))
		newPassword := req.PostFormValue("new_password")
		if newPassword != "" {
			if err = authService.SetCredential(user.ID, &Credential{
				Type:  CredentialTypePassword,
				Value: newPassword,
			}); err != nil {
				log.WithField("Error", err).Warning("Failed to update password")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if email != "" {
			update := &User{
				ID:    user.ID,
				Email: email,
			}

			if err = userService.UpdateUser(user, update); err != nil {
				log.WithField("Error", err).Warning("Failed to update user")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			user = update
		}

		// create new token
		authToken, err := authService.NewToken(user)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to generate token")
			res.WriteHeader(http.StatusInternalServerError)
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

		response.Respond(
			res,
			req,
			response.NewRedirectResponse(publicPath, strings.TrimSuffix(req.URL.RequestURI(), "/edit")),
		)
	}
}
