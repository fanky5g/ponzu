// Package editor enables users to create edit templates from their entities
// structs so that admins can manage entities
package editor

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/fanky5g/ponzu/internal/templates"
)

// Editable ensures data is editable
type Editable interface {
	MarshalEditor(publicPath string) ([]byte, error)
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
	Parent                 string
	TypeName               string
	PositionalPlaceHolders []string
}

func makeScript(name string) *template.Template {
	t, err := templates.Template(fmt.Sprintf("scripts/%s", name))
	if err != nil {
		panic(err)
	}

	return t
}
