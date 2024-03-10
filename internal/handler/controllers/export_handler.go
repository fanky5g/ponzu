package controllers

import (
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/handler/controllers/views"
	"github.com/fanky5g/ponzu/internal/services/config"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/internal/util"
	"io"
	"log"
	"net/http"
	"strings"
)

func NewExportHandler(pathConf conf.Paths, configService config.Service, contentService content.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		appName, err := configService.GetAppName()
		if err != nil {
			log.Printf("Failed to get app name: %v\n", appName)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// /contents/export?type=Blogpost&format=csv
		q := req.URL.Query()
		t := q.Get("type")
		f := strings.ToLower(q.Get("format"))

		if t == "" || f == "" {
			v, err := views.Admin(util.Html("error_400"), appName, pathConf)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			res.WriteHeader(http.StatusBadRequest)
			_, err = res.Write(v)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

		}

		switch f {
		case "csv":
			response, err := contentService.ExportCSV(t)
			if err != nil {
				LogAndFail(res, err, appName, pathConf)
			}

			if response == nil {
				res.WriteHeader(http.StatusNoContent)
				return
			}

			res.Header().Set("Content-Type", response.ContentType)
			res.Header().Set("Content-Disposition", response.ContentDisposition)
			if _, err = io.Copy(res, response.Payload); err != nil {
				LogAndFail(res, err, appName, pathConf)
			}
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
