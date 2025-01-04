package response

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	data interface{}
	err  error
}

func (r *jsonResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	payload := map[string]interface{}{
		"data": r.data,
	}

	if r.err != nil {
		payload = map[string]interface{}{
			"error": map[string]string{
				"message": r.err.Error(),
			},
		}
	}

	return json.NewEncoder(w).Encode(payload)
}

func NewJSONResponse(statusCode int, data interface{}, err error) *Response {
	return &Response{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Vary":         "Accept-Encoding",
		},
		Renderer: &jsonResponse{
			data: data,
			err:  err,
		},
	}
}
