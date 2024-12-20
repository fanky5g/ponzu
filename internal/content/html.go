package content

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func WriteTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	w.Header().Set("Content-Type", "text/Html")
	if err := tmpl.Execute(w, data); err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Failed to write response")
		return
	}
}
