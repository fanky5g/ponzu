package auth

import (
	"github.com/fanky5g/ponzu/internal/http/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// POST
func NewAPIAuthHandler(authService *Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		accountId, credential, err := MapAuthRequest(req)
		if err != nil {
			log.WithField("err", err).Error("error occurred during authentication")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		authToken, err := authService.LoginByEmail(accountId, credential)
		if err != nil {
			response.Respond(
				res,
				req,
				response.NewJSONResponse(http.StatusBadRequest, nil, err),
			)
			return
		}

		response.Respond(
			res,
			req,
			response.NewJSONResponse(http.StatusOK, MapAuthTokenResponse(authToken), nil),
		)
	}
}
