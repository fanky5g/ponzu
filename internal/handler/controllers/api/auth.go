package api

import (
	"github.com/fanky5g/ponzu/internal/application/auth"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/response"

	"net/http"
)

func NewAuthHandler(authService auth.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			accountId, credential, err := request.MapAuthRequest(req)
			if err != nil {
				writeJSONError(res, http.StatusBadRequest, err)
				return
			}

			authToken, err := authService.LoginByEmail(accountId, credential)
			if err != nil {
				writeJSONError(res, http.StatusBadRequest, err)
				return
			}

			writeJSONData(res, http.StatusOK, response.MapAuthTokenResponse(authToken))
			return
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
