package middleware

import (
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/handler/controllers/mappers/request"
	"github.com/fanky5g/ponzu/internal/services/auth"
	"github.com/fanky5g/ponzu/internal/util"
	"log"
	"net/http"
)

var AuthMiddleware Token = "AuthMiddleware"

func NewAuthMiddleware(paths conf.Paths, authService auth.Service) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			authToken := request.GetAuthToken(req)
			isValid, err := authService.IsTokenValid(authToken)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				log.Printf("Failed to check token validity: %v\n", err)
				return
			}

			if isValid {
				next.ServeHTTP(res, req)
				return
			}

			util.Redirect(req, res, paths, "/login", http.StatusFound)
		}
	}
}
