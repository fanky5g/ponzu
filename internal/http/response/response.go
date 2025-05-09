package response

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ResponseRenderer interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

type Response struct {
	Headers    map[string]string
	StatusCode int
	Renderer   ResponseRenderer
}

func Respond(w http.ResponseWriter, r *http.Request, response *Response) {
	if response.StatusCode != 0 && response.StatusCode != http.StatusOK {
		w.WriteHeader(response.StatusCode)
	}

	for headerName, headerValue := range response.Headers {
		w.Header().Set(headerName, headerValue)
	}

	// w.Header().Add("Content-Type", response.ContentType)
	// w.Header().Add("Content-Disposition", response.ContentDisposition)

	if err := response.Renderer.Render(w, r); err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Failed to write response")
		return
	}
}
