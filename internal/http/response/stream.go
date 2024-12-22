package response

import (
	"io"
	"net/http"

	"github.com/fanky5g/ponzu/internal/datasource"
)

type stream struct {
	data io.Reader
}

func (s *stream) Render(w http.ResponseWriter, r *http.Request) error {
	var err error
	_, err = io.Copy(w, s.data)

	return err
}

func NewStreamResponse(statusCode int, ds datasource.Datasource) *Response {
	return &Response{
		StatusCode:         statusCode,
		ContentType:        ds.GetContentType(),
		ContentDisposition: ds.GetContentDisposition(),
		Renderer: &stream{
			data: ds,
		},
	}
}
