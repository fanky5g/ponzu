package util

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	pathToTemplates string
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../..")
	pathToTemplates = fmt.Sprintf("%s/internal/handler/controllers/views", rootPath)
}

func MakeTemplate(templates ...string) *template.Template {
	return template.Must(template.New(strings.Join(templates, "_")).Parse(Html(templates...)))
}

func MakeScript(name string) *template.Template {
	templateString, err := getTemplateString(name, "scripts")
	if err != nil {
		panic(err)
	}

	return template.Must(template.New(name).Parse(templateString))
}

func Html(templates ...string) string {
	var tmpl string
	for _, name := range templates {
		htmlString, err := getTemplateString(name, "html")
		if err != nil {
			panic(err)
		}

		tmpl += htmlString
	}

	return tmpl
}

func getTemplateString(name, templateType string) (string, error) {
	templateName := fmt.Sprintf("%s/%s/%s.gohtml", pathToTemplates, templateType, name)
	f, err := os.Open(templateName)
	if err != nil {
		return "", fmt.Errorf("failed to open template file: %s. Error: %v", name, err)
	}

	htmlBytes, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	return string(htmlBytes), err
}
