package controllers

import (
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"net/http"
	"time"
)

func NewLogoutHandler(r router.Router) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:    "_token",
			Expires: time.Unix(0, 0),
			Value:   "",
			Path:    "/",
		})

		r.Redirect(req, res, "/login")
	}
}
