package controllers

import (
	"bytes"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/analytics"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewAdminHandler(r router.Router) http.HandlerFunc {
	analyticsService := r.Context().Service(tokens.AnalyticsServiceToken).(analytics.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		buf := &bytes.Buffer{}
		data, err := analyticsService.GetChartData()
		if err != nil {
			log.WithFields(log.Fields{"Error": err}).Warn("Failed to get chart data")
			r.Renderer().InternalServerError(res)
			return
		}

		tmpl := r.Renderer().Template("analytics")
		err = tmpl.Execute(buf, data)
		if err != nil {
			log.WithFields(log.Fields{"Error": err}).Warn("Failed to make template")
			r.Renderer().InternalServerError(res)
			return
		}

		r.Renderer().InjectInAdminView(res, buf)
	}
}
