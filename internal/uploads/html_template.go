package uploads

import (
	"html/template"
	"path/filepath"
	"runtime"

	"github.com/fanky5g/ponzu/internal/views"
)

func getEditUploadTemplate() (*template.Template, error) {
	_, b, _, _ := runtime.Caller(0)
	sharedTemplatesRoot := filepath.Join(filepath.Dir(b), "../dashboard")

	return template.New("edit-upload").Funcs(views.GlobFuncs).Parse(
		views.Html(
			filepath.Join(sharedTemplatesRoot, "dashboard.gohtml"),
			filepath.Join(sharedTemplatesRoot, "app-frame.gohtml"),
			filepath.Join(filepath.Dir(b), "edit_upload_view.gohtml"),
		),
	)
}
