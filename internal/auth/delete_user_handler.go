package auth

import (
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/http/response"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func NewDeleteUserHandler(publicPath string, authService *Service, userService *UserService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse form")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// do not allow current user to delete themselves
		user, err := authService.GetUserFromAuthToken(request.GetAuthToken(req))
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get auth token")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		email := strings.ToLower(req.PostFormValue("email"))
		if user.Email == email {
			log.Debug("cannot delete own user account")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		// delete existing user
		err = userService.DeleteUser(email)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to delete user")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewRedirectResponse(
				publicPath,
				strings.TrimSuffix(req.URL.RequestURI(), "/delete"),
			))
	}
}
