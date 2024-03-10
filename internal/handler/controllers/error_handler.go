package controllers

import (
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/handler/controllers/views"
	"github.com/fanky5g/ponzu/internal/util"
	"log"
	"net/http"
)

func writeResponse(res http.ResponseWriter, statusCode int, response []byte) {
	res.WriteHeader(statusCode)
	if _, err := res.Write(response); err != nil {
		log.Printf("Error writing response: %v\n", err)
	}
}

func renderErrorView(res http.ResponseWriter, appName, templateName string, statusCode int, pathConf conf.Paths) {
	errView, err := views.Admin(util.Html(templateName), appName, pathConf)
	if err != nil {
		log.Printf("Failed to build error view: %v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeResponse(res, statusCode, errView)
	return
}

func LogAndFail(res http.ResponseWriter, err error, appName string, pathConf conf.Paths) {
	if err != nil {
		log.Println(err)
		renderErrorView(res, appName, "error_500", http.StatusInternalServerError, pathConf)
	}
}
