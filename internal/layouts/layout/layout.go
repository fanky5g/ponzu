package layout

import (
	"github.com/fanky5g/ponzu/internal/templates"
	"html/template"
	"io"
)

type TemplateExecutor func(tmpl *template.Template, w io.Writer, data interface{}) error

type Layout struct {
	tmpl *template.Template
	exec TemplateExecutor
}

func (layout *Layout) Child(names ...string) (*Layout, error) {
	c, err := layout.tmpl.Clone()
	if err != nil {
		return nil, err
	}

	t, err := c.Parse(templates.Html(names...))
	if err != nil {
		return nil, err
	}

	return &Layout{
		tmpl: t,
		exec: layout.exec,
	}, nil
}

func (layout *Layout) Execute(w io.Writer, data interface{}) error {
	return layout.exec(layout.tmpl, w, data)
}

func NewLayout(tmpl *template.Template, exec TemplateExecutor) *Layout {
	return &Layout{
		tmpl: tmpl,
		exec: exec,
	}
}
