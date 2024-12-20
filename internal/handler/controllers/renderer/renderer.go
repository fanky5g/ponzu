package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/manager"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/viewparams/table"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router/context"
	"github.com/fanky5g/ponzu/internal/services/config"
	"github.com/fanky5g/ponzu/internal/views"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
)

type View struct {
	Logo       string
	PublicPath string
	Types      map[string]content.Builder
	Subview    template.HTML
	Data       interface{}
}

type Renderer interface {
	InjectInAdminView(w http.ResponseWriter, subView *bytes.Buffer)
	Render(w http.ResponseWriter, view string)
	InjectTemplateInAdmin(res http.ResponseWriter, templateText string, data interface{})
	Editable(http.ResponseWriter, editor.Editable)
	InternalServerError(w http.ResponseWriter)
	MethodNotAllowed(w http.ResponseWriter)
	BadRequest(res http.ResponseWriter)
	ManageEditable(res http.ResponseWriter, editable editor.Editable, typeName string)
	Html(res http.ResponseWriter, data []byte)
	Json(res http.ResponseWriter, statusCode int, data interface{})
	Error(res http.ResponseWriter, statusCode int, err error)
	Template(templates ...string) *template.Template
	TemplateString(templates ...string) string
	TemplateFromDir(name string) *template.Template
	TableView(w http.ResponseWriter, templateName string, params *table.Table)
}

type renderer struct {
	ctx context.Context
}

func (r *renderer) InjectInAdminView(res http.ResponseWriter, subView *bytes.Buffer) {
	configService := r.ctx.Service(tokens.ConfigServiceToken).(config.Service)
	appName, err := configService.GetAppName()

	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Failed to get app name")
		r.InternalServerError(res)
		return
	}

	a := View{
		Logo:       appName,
		Types:      r.ctx.Types().Content,
		Subview:    template.HTML(subView.Bytes()),
		PublicPath: r.ctx.Paths().PublicPath,
	}

	buf := &bytes.Buffer{}
	tmpl := r.Template("start_admin", "main_admin", "end_admin")
	err = tmpl.Execute(buf, a)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Failed to execute template")
		r.InternalServerError(res)
		return
	}

	r.Html(res, buf.Bytes())
}

func (r *renderer) Render(res http.ResponseWriter, templateName string) {
	view, err := r.renderInAppFrame(templateName)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Failed to render view")
		r.InternalServerError(res)
		return
	}

	r.Html(res, view)
}

func (r *renderer) InternalServerError(res http.ResponseWriter) {
	errView, err := r.renderInAppFrame("error_500")
	if err != nil {
		log.Printf("Failed to build error view: %v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusInternalServerError)
	if _, err = res.Write(errView); err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Error writing response")
	}
}

func (r *renderer) BadRequest(res http.ResponseWriter) {
	errView, err := r.renderInAppFrame("error_400")
	if err != nil {
		log.Printf("Failed to build error view: %v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusBadRequest)
	if _, err = res.Write(errView); err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Error writing response")
	}
}

func (r *renderer) MethodNotAllowed(res http.ResponseWriter) {
	errView, err := r.renderInAppFrame("error_405")
	if err != nil {
		log.Printf("Failed to build error view: %v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusMethodNotAllowed)
	if _, err = res.Write(errView); err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Error writing response")
	}
}

func (r *renderer) Editable(res http.ResponseWriter, editable editor.Editable) {
	b, err := editable.MarshalEditor(r.ctx.Paths().PublicPath)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Failed to Render editable")
		r.InternalServerError(res)
		return
	}

	r.InjectInAdminView(res, bytes.NewBuffer(b))
}

func (r *renderer) ManageEditable(res http.ResponseWriter, editable editor.Editable, typeName string) {
	m, err := manager.Manage(editable, r.ctx.Paths().PublicPath, typeName)
	if err != nil {
		log.WithField("Error", err).Warning("Failed to execute editable manager")
		return
	}

	r.InjectInAdminView(res, bytes.NewBuffer(m))
}

func (r *renderer) InjectTemplateInAdmin(res http.ResponseWriter, templateText string, data interface{}) {
	view := View{
		PublicPath: r.ctx.Paths().PublicPath,
		Data:       data,
	}

	tmpl, err := template.New("template").Parse(templateText)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Failed to create template")
		r.InternalServerError(res)
		return
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, view)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Warn("Failed to execute template")
		r.InternalServerError(res)
		return
	}

	r.InjectInAdminView(res, buf)
}

func (r *renderer) renderInAppFrame(template string) ([]byte, error) {
	configService := r.ctx.Service(tokens.ConfigServiceToken).(config.Service)
	appName, err := configService.GetAppName()
	if err != nil {
		return nil, err
	}

	view := View{
		Logo:       appName,
		Types:      r.ctx.Types().Content,
		PublicPath: r.ctx.Paths().PublicPath,
	}

	buf := &bytes.Buffer{}
	tmpl := r.Template("start_admin", template, "end_admin")
	err = tmpl.Execute(buf, view)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *renderer) Template(templates ...string) *template.Template {
	templatePaths := make([]string, len(templates))
	for i, templateName := range templates {
		templatePaths[i] = fmt.Sprintf("%s/%s", views.Path, templateName)
	}

	return template.Must(template.New(strings.Join(templates, "_")).
		Funcs(views.GlobFuncs).
		Parse(views.Html(templatePaths...)))
}

func (r *renderer) TemplateString(templates ...string) string {
	templatePaths := make([]string, len(templates))
	for i, templateName := range templates {
		templatePaths[i] = fmt.Sprintf("%s/%s", views.Path, templateName)
	}

	return views.Html(templatePaths...)
}

// Dir recursively walks directory and returns all go template files.
func (r *renderer) TemplateFromDir(name string) *template.Template {
	goTemplates := make([]string, 0)
	rootPath := filepath.Join(views.Path, name)
	if err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".gohtml") {
			goTemplates = append(goTemplates, filepath.Join(name, d.Name()))
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return r.Template(goTemplates...)
}

func New(ctx context.Context) (Renderer, error) {
	return &renderer{
		ctx: ctx,
	}, nil
}
