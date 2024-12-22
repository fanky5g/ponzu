package response

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ResponseRenderer interface {
	Render(w http.ResponseWriter) error
}

type Response struct {
	ContentType        string
	ContentDisposition string
	StatusCode         int
	Renderer           ResponseRenderer
}

func Write(w http.ResponseWriter, response *Response) {
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(response.StatusCode)
	}

	w.Header().Add("Content-Type", response.ContentType)
	w.Header().Add("Content-Disposition", response.ContentDisposition)

	if err := response.Renderer.Render(w); err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Failed to write response")
		return
	}
}
