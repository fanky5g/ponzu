// Package editor enables users to create edit templates from their entities
// structs so that admins can manage entities
package editor

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/util"
	"html/template"
	"path/filepath"
	"runtime"
)

var pathToTemplates string

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../..")
	pathToTemplates = fmt.Sprintf("%s/content/editor/templates", rootPath)
}

// Editable ensures data is editable
type Editable interface {
	MarshalEditor(paths config.Paths) ([]byte, error)
}

// Editor is a view containing fields to manage entities
type Editor struct {
	ViewBuf *bytes.Buffer
}

type ContentMetadata struct {
	TypeName string
}

// Field is used to create the editable view for a field
// within a particular entities struct
type Field struct {
	View []byte
}

type FieldArgs struct {
	Parent   string
	TypeName string
}

func makeScript(name string) *template.Template {
	templateString := util.Html(fmt.Sprintf("%s/scripts/%s", pathToTemplates, name))

	return template.Must(template.New(name).Parse(templateString))
}

func makeHtml(name string) string {
	return util.Html(fmt.Sprintf("%s/html/%s", pathToTemplates, name))
}
