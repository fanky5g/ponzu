package content

import (
	"html/template"
	"maps"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/views"
)

var (
	pageTemplates []string
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	sharedTemplatesRoot := filepath.Join(filepath.Dir(b), "../dashboard")

	pageTemplates = []string{
		filepath.Join(sharedTemplatesRoot, "dashboard.gohtml"),
		filepath.Join(sharedTemplatesRoot, "app-frame.gohtml"),
		filepath.Join(filepath.Dir(b), "edit_content_view.gohtml"),
	}
}

func NewEditContentFormHandler(propCache config.ApplicationPropertiesCache, contentService *Service) http.HandlerFunc {
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
		contentQuery, err := MapToContentQuery(req)
		if err != nil {
			// TODO: handle error
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		typeBuilder, ok := contentTypes[contentQuery.Type]
		if !ok {
			// TODO: handle error
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		entity := typeBuilder()
		if contentQuery.ID != "" {
			entity, err = contentService.GetContent(contentQuery.Type, contentQuery.ID)
			if err != nil {
				// TODO: handle error
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		editContentForm, err := NewEditContentFormViewModel(entity, propCache)
		if err != nil {
			// TODO: handle error
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		WriteTemplate(res, editPage, editContentForm)
		return
	}
}
