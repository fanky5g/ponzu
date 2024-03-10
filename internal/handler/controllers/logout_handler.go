package controllers

import (
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/util"
	"net/http"
	"time"
)

func NewLogoutHandler(pathConf conf.Paths) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:    "_token",
			Expires: time.Unix(0, 0),
			Value:   "",
			Path:    "/",
		})

		util.Redirect(req, res, pathConf, "/login", http.StatusFound)
	}
}
