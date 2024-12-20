package edit

import (
	"html/template"
	"maps"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/fanky5g/ponzu/internal/handler/controllers/dashboard/mapper"
	"github.com/fanky5g/ponzu/internal/handler/controllers/dashboard/renderer"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/internal/views"
)

var (
	pageTemplates []string
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	sharedTemplatesRoot := filepath.Join(filepath.Dir(b), "../../template")

	pageTemplates = []string{
		filepath.Join(sharedTemplatesRoot, "admin.gohtml"),
		filepath.Join(sharedTemplatesRoot, "app-frame.gohtml"),
		filepath.Join(filepath.Dir(b), "edit_content_form.gohtml"),
	}
}

func NewEditContentFormHandler(
	propCache config.ApplicationPropertiesCache,
	contentService content.Service,
) http.HandlerFunc {
	funcs := views.GlobFuncs
	maps.Copy(funcs, workflow.TemplateFuncs)

	editPage := template.Must(
		template.New("edit").Funcs(funcs).Parse(views.Html(pageTemplates...)),
	)

	contentTypes, err := propCache.GetContentTypes()
	if err != nil {
		panic(err)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		identifier, err := mapper.MapToContentIdentifier(req)
		if err != nil {
			// TODO: handle error
			return
		}

		typeBuilder, ok := contentTypes[identifier.Type]
		if !ok {
			// TODO: handle error
			return
		}

		entity := typeBuilder()
		if identifier.ID != "" {
			entity, err = contentService.GetContent(identifier.Type, identifier.ID)
			if err != nil {
				// TODO: handle error
				return
			}
		}

		editContentForm, err := NewEditContentForm(entity, propCache)
		if err != nil {
			// TODO: handle error
			return
		}

		renderer.WriteTemplate(res, editPage, editContentForm)
		return
	}
}
