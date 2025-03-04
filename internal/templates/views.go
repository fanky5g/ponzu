package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"strings"
)

var (
	//go:embed all:*
	templates embed.FS
)

func ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	t, err := Template(name)
	if err != nil {
		return err
	}

	return t.Execute(w, data)
}

func Template(name string) (*template.Template, error) {
	templateData, err := getTemplateString(name)
	if err != nil {
		return nil, err
	}

	return template.New(name).Funcs(GlobFuncs).Parse(templateData)
}

func Glob(name, pattern string) (*template.Template, error) {
	return template.New(name).Funcs(GlobFuncs).ParseFS(templates, pattern)
}

func Html(names ...string) string {
	var tmpl string
	for _, name := range names {
		htmlString, err := getTemplateString(name)
		if err != nil {
			panic(err)
		}

		tmpl += htmlString
	}

	return tmpl
}

func getTemplateString(path string) (string, error) {
	normalizedPath := fmt.Sprintf("%s%s", strings.TrimSuffix(path, ".gohtml"), ".gohtml")
	f, err := templates.Open(normalizedPath)
	if err != nil {
		return "", fmt.Errorf("failed to open template file: %s. Error: %v", path, err)
	}

	htmlBytes, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	return string(htmlBytes), err
}
