package middleware

import (
	"net/http"

	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/services"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/fanky5g/ponzu/util"
	log "github.com/sirupsen/logrus"
)

var AuthMiddleware Token = "AuthMiddleware"

func NewAuthMiddleware(paths conf.Paths, applicationServices services.Services) func(next http.HandlerFunc) http.HandlerFunc {
	authService := applicationServices.Get(tokens.AuthServiceToken).(auth.Service)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			authToken := request.GetAuthToken(req)
			isValid, err := authService.IsTokenValid(authToken)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				log.WithField("Error", err).Warning("Failed to check token validity")
				return
			}

			if isValid {
				next.ServeHTTP(res, req)
				return
			}

			routeTag, hasRouteTag := req.Context().Value(constants.RouteTagIdentifier).(constants.RouteTag)
			if hasRouteTag && routeTag == constants.APIRoute {
				util.WriteJSONResponse(res, http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]string{
						"message": "Unauthorized",
					},
				})

				return
			}

			util.Redirect(req, res, paths, "/login", http.StatusFound)
		}
	}
}
