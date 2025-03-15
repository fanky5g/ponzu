package root

import (
	"github.com/fanky5g/ponzu/internal/layouts/layout"
	"github.com/fanky5g/ponzu/internal/templates"
	"html/template"
	"io"
)

func GetTemplate() (*template.Template, error) {
	return template.New("root").Funcs(templates.GlobFuncs).Parse(
		templates.Html("views/root.gohtml"),
	)
}

type AppNameProvider interface {
	GetAppName() (string, error)
}

func New(appNameProvider AppNameProvider, publicPath string) (*layout.Layout, error) {
	t, templateErr := GetTemplate()
	if templateErr != nil {
		return nil, templateErr
	}

	return layout.NewLayout(t, func(tmpl *template.Template, w io.Writer, data interface{}) error {
		viewModel, err := NewRootViewModel(appNameProvider, publicPath)
		if err != nil {
			return err
		}

		return tmpl.Execute(w, struct {
			*ViewModel
			Data interface{}
		}{
			ViewModel: viewModel,
			Data:      data,
		})
	}), nil
}
