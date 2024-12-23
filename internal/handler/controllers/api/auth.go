package api

import (
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/tokens"
	"net/http"
)

func NewAuthHandler(r router.Router) http.HandlerFunc {
	authService := r.Context().Service(tokens.AuthServiceToken).(auth.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			accountId, credential, err := request.MapAuthRequest(req)
			if err != nil {
				r.Renderer().Error(res, http.StatusBadRequest, err)
				return
			}

			authToken, err := authService.LoginByEmail(accountId, credential)
			if err != nil {
				r.Renderer().Error(res, http.StatusBadRequest, err)
				return
			}

			r.Renderer().Json(res, http.StatusOK, response.MapAuthTokenResponse(authToken))
			return
		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
