package content

import (
	"github.com/fanky5g/ponzu/internal/templates"
	"html/template"
	"sync"
)

var (
	editUploadTemplate        *template.Template
	editUploadTemplateCreator sync.Once
)

func getEditUploadTemplate(layoutTmpl *template.Template) (*template.Template, error) {
	var err error

	editUploadTemplateCreator.Do(func() {
		editUploadTemplate, err = layoutTmpl.Parse(templates.Html("views/edit_upload_view.gohtml"))
	})

	return editUploadTemplate, err
}

// TODO(B.B) use dashboard layout template
func getEditPageTemplate() (*template.Template, error) {
	return template.New("edit").Funcs(templates.GlobFuncs).Parse(
		templates.Html(
			"views/dashboard.gohtml",
			"views/app-frame.gohtml",
			"views/edit_content_view.gohtml",
		),
	)
}
