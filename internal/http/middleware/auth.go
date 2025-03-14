package middleware

import (
	"github.com/fanky5g/ponzu/internal/http/response"
	"net/http"

	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/util"
	log "github.com/sirupsen/logrus"
)

var AuthMiddleware Token = "AuthMiddleware"

type TokenValidatorInterface interface {
	IsTokenValid(authToken string) (bool, error)
}

func NewAuthMiddleware(paths conf.Paths, tokenValidator TokenValidatorInterface) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			authToken := request.GetAuthToken(req)
			isValid, err := tokenValidator.IsTokenValid(authToken)
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

			response.Respond(res, req, response.NewRedirectResponse(paths.PublicPath, "/login"))
		}
	}
}
