package views

import (
	"bytes"
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/util"
	"html/template"
)

type View struct {
	Logo       string
	PublicPath string
	Types      map[string]item.EntityBuilder
	Subview    template.HTML
}

// Admin ...
func Admin(view, appName string, paths conf.Paths) (_ []byte, err error) {
	a := View{
		Logo:       appName,
		Types:      item.Types,
		Subview:    template.HTML(view),
		PublicPath: paths.PublicPath,
	}

	buf := &bytes.Buffer{}
	tmpl := util.MakeTemplate("start_admin", "main_admin", "end_admin")
	err = tmpl.Execute(buf, a)
	if err != nil {
		return
	}

	return buf.Bytes(), nil
}
