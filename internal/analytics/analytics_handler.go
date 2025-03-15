package analytics

import (
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewAnalyticsHandler(analyticsService *Service, layout layouts.Template) http.HandlerFunc {
	tmpl, templateErr := layout.Child("views/analytics")
	if templateErr != nil {
		panic(templateErr)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		data, err := analyticsService.GetChartData()
		if err != nil {
			log.WithFields(log.Fields{"Error": err}).Warn("Failed to get chart data")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(res, req, response.NewHTMLResponse(http.StatusOK, tmpl, data))
	}
}
