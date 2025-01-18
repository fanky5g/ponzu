package content

import (
	"github.com/fanky5g/ponzu/internal/views"
	"html/template"
	"maps"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	editUploadTemplate        *template.Template
	editUploadTemplateCreator sync.Once
)

func getEditUploadTemplate(layoutTmpl *template.Template) (*template.Template, error) {
	var err error

	editUploadTemplateCreator.Do(func() {
		_, b, _, _ := runtime.Caller(0)
		workingDirectory := filepath.Dir(b)

		editUploadTemplate, err = layoutTmpl.Parse(views.Html(filepath.Join(workingDirectory, "templates/edit_upload_view.gohtml")))
	})

	return editUploadTemplate, err
}

// TODO(B.B) use dashboard layout template
func getEditPageTemplate() (*template.Template, error) {
	_, b, _, _ := runtime.Caller(0)
	workingDirectory := filepath.Dir(b)
	sharedTemplatesRoot := filepath.Join(workingDirectory, "../dashboard")

	funcs := views.GlobFuncs
	maps.Copy(funcs, TemplateFuncs)

	return template.New("edit").Funcs(funcs).Parse(
		views.Html(
			filepath.Join(sharedTemplatesRoot, "dashboard.gohtml"),
			filepath.Join(sharedTemplatesRoot, "app-frame.gohtml"),
			filepath.Join(workingDirectory, "templates/edit_content_view.gohtml"),
		),
	)
}
