package renderer

import (
	"github.com/fanky5g/ponzu/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (r *renderer) Html(res http.ResponseWriter, data []byte) {
	res.Header().Set("Content-Type", "text/Html")
	if _, err := res.Write(data); err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Failed to write response")
		r.InternalServerError(res)
		return
	}
}

func (r *renderer) Json(res http.ResponseWriter, statusCode int, data interface{}) {
	util.WriteJSONResponse(res, statusCode, map[string]interface{}{
		"data": data,
	})
}

func (r *renderer) Error(res http.ResponseWriter, statusCode int, err error) {
	util.WriteJSONResponse(res, statusCode, map[string]interface{}{
		"error": map[string]string{
			"message": err.Error(),
		},
	})
}
