package content

import (
	"html/template"
	"maps"
	"path/filepath"
	"runtime"

	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/fanky5g/ponzu/internal/views"
)

func getEditPageTemplate() (*template.Template, error) {
	_, b, _, _ := runtime.Caller(0)
	sharedTemplatesRoot := filepath.Join(filepath.Dir(b), "../dashboard")

	funcs := views.GlobFuncs
	maps.Copy(funcs, workflow.TemplateFuncs)

	return template.New("edit").Funcs(funcs).Parse(
		views.Html(
			filepath.Join(sharedTemplatesRoot, "dashboard.gohtml"),
			filepath.Join(sharedTemplatesRoot, "app-frame.gohtml"),
			filepath.Join(filepath.Dir(b), "edit_content_view.gohtml"),
		),
	)
}
