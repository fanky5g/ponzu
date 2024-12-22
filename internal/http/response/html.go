package response

import (
	"net/http"
	"html/template"
)

type templateRenderer struct {
	data interface{}
	tmpl *template.Template
}

func (renderer *templateRenderer) Render(w http.ResponseWriter) error {
	return renderer.tmpl.Execute(w, renderer.data)
}

func NewHTMLResponse(statusCode int, tmpl *template.Template, data interface{}) *Response {
	return &Response{
		StatusCode:  statusCode,
		ContentType: "text/html",
		Renderer: &templateRenderer{
			tmpl: tmpl,
			data: data,
		},
	} 
}
