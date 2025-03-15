package response

import (
	"io"
	"net/http"
)

type TemplateExecutorInterface interface {
	Execute(w io.Writer, data interface{}) error
}

type templateRenderer struct {
	data interface{}
	tmpl TemplateExecutorInterface
}

func (renderer *templateRenderer) Render(w http.ResponseWriter, r *http.Request) error {
	return renderer.tmpl.Execute(w, renderer.data)
}

func NewHTMLResponse(statusCode int, tmpl TemplateExecutorInterface, data interface{}) *Response {
	return &Response{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Renderer: &templateRenderer{
			tmpl: tmpl,
			data: data,
		},
	}
}
