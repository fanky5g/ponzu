package content

import (
	"github.com/fanky5g/ponzu/internal/views"
	"html/template"
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

		editUploadTemplate, err = layoutTmpl.Parse(views.Html(filepath.Join(workingDirectory, "edit_upload_view.gohtml")))
	})

	return editUploadTemplate, err
}
