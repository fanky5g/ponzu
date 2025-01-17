package uploads

import (
	"github.com/fanky5g/ponzu/internal/views"
	"html/template"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	tmpl *template.Template
	once sync.Once
)

func getEditUploadTemplate(layoutTmpl *template.Template) (*template.Template, error) {
	var err error

	once.Do(func() {
		_, b, _, _ := runtime.Caller(0)
		workingDirectory := filepath.Dir(b)

		tmpl, err = layoutTmpl.Parse(views.Html(filepath.Join(workingDirectory, "edit_upload_view.gohtml")))
	})

	return tmpl, err
}
