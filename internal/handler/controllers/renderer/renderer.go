package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/viewparams/table"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router/context"
	"github.com/fanky5g/ponzu/internal/templates"
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
	configCache := r.ctx.Service(tokens.ConfigCache).(config.ConfigCache)
	appName, err := configCache.GetAppName()
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
	configCache := r.ctx.Service(tokens.ConfigCache).(config.ConfigCache)
	appName, err := configCache.GetAppName()
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

func (r *renderer) Template(names ...string) *template.Template {
	templatePaths := make([]string, len(names))
	for i, templateName := range names {
		templatePaths[i] = fmt.Sprintf("views/%s", templateName)
	}

	return template.Must(template.New(strings.Join(names, "_")).
		Funcs(templates.GlobFuncs).
		Parse(templates.Html(templatePaths...)))
}

func (r *renderer) TemplateString(names ...string) string {
	templatePaths := make([]string, len(names))
	for i, templateName := range names {
		templatePaths[i] = fmt.Sprintf("views/%s", templateName)
	}

	return templates.Html(templatePaths...)
}

// TemplateFromDir recursively walks directory and returns all go template files.
func (r *renderer) TemplateFromDir(dirName string) *template.Template {
	pattern := fmt.Sprintf("views/%s/*.gohtml", dirName)
	t, err := templates.Glob(dirName, pattern)
	if err != nil {
		panic(err)
	}

	return t
}

func New(ctx context.Context) (Renderer, error) {
	return &renderer{
		ctx: ctx,
	}, nil
}
