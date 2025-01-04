package response

import (
	"html/template"
	"net/http"
)

type templateRenderer struct {
	data interface{}
	tmpl *template.Template
}

func (renderer *templateRenderer) Render(w http.ResponseWriter, r *http.Request) error {
	return renderer.tmpl.Execute(w, renderer.data)
}

func NewHTMLResponse(statusCode int, tmpl *template.Template, data interface{}) *Response {
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
