package views

import (
	"embed"
	"html/template"
	"io"
	"path/filepath"
	"runtime"
)

var (
	t *template.Template
	//go:embed all:*.gohtml
	templates embed.FS
	Path      string
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	Path = filepath.Join(filepath.Dir(b), "../../internal/views")

	var err error
	t, err = template.New("views").Funcs(GlobFuncs).ParseFS(templates, "*.gohtml")
	if err != nil {
		panic(err)
	}
}

func ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	return t.ExecuteTemplate(w, name, data)
}
