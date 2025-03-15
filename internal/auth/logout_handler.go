package auth

import (
	"github.com/fanky5g/ponzu/internal/http/response"
	"net/http"
	"time"
)

func NewLogoutHandler(publicPath string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:    "_token",
			Expires: time.Unix(0, 0),
			Value:   "",
			Path:    "/",
		})

		response.Respond(res, req, response.NewRedirectResponse(publicPath, "/login"))
	}
}
