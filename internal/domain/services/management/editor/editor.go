// Package editor enables users to create edit views from their content
// structs so that admins can manage content
package editor

import (
	"bytes"
	"github.com/fanky5g/ponzu/internal/util"
	"log"
)

// Editable ensures data is editable
type Editable interface {
	MarshalEditor() ([]byte, error)
}

// Editor is a view containing fields to manage content
type Editor struct {
	ViewBuf *bytes.Buffer
}

// Field is used to create the editable view for a field
// within a particular content struct
type Field struct {
	View []byte
}

type FieldArgs struct {
	Parent   string
	TypeName string
}

// Form takes editable content and any number of Field funcs to describe the edit
// page for any content struct added by a user
func Form(post Editable, fields ...Field) ([]byte, error) {
	editor := &Editor{}

	editor.ViewBuf = &bytes.Buffer{}
	_, err := editor.ViewBuf.WriteString(`
		<div class="row">
			<div class="col s12">
			  <ul class="tabs">
				<li class="tab col s3"><a class="active" href="#content"><i class="material-icons">edit</i>Edit</a></li>
				<li class="tab col s3"><a href="#properties"><i class="material-icons">tune</i>Properties</a></li>
			  </ul>
			</div>
			<div id="content" class="col s12 editor-content">
	`)

	if err != nil {
		log.Println("Error writing HTML string to editor Form buffer")
		return nil, err
	}

	for _, f := range fields {
		if err = addFieldToEditorView(editor, f); err != nil {
			return nil, err
		}
	}

	_, err = editor.ViewBuf.WriteString(`</div>`)
	if err != nil {
		log.Println("Error writing HTML string to editor Form buffer")
		return nil, err
	}

	// content items with Item embedded have some default fields we need to render
	_, err = editor.ViewBuf.WriteString(`<div id="properties" class="col s12 editor-metadata">`)
	if err != nil {
		log.Println("Error writing HTML string to editor Form buffer")
		return nil, err
	}

	contentMetadata := util.Html("editor_content_metadata")

	_, err = editor.ViewBuf.WriteString(contentMetadata)
	if err != nil {
		log.Println("Error writing HTML string to editor Form buffer")
		return nil, err
	}

	err = addPostDefaultFieldsToEditorView(post, editor)
	if err != nil {
		return nil, err
	}

	_, err = editor.ViewBuf.WriteString(`</div>`)
	if err != nil {
		log.Println("Error writing HTML string to editor Form buffer")
		return nil, err
	}

	script := &bytes.Buffer{}
	scriptTmpl := util.MakeScript("editor")
	if err = scriptTmpl.Execute(script, nil); err != nil {
		panic(err)
	}

	editorControls := util.Html("editor_controls")
	_, err = editor.ViewBuf.WriteString(editorControls + script.String() + `</div>`)

	return editor.ViewBuf.Bytes(), nil
}

func addFieldToEditorView(e *Editor, f Field) error {
	_, err := e.ViewBuf.Write(f.View)
	if err != nil {
		log.Println("Error writing field view to editor view buffer")
		return err
	}

	return nil
}

func addPostDefaultFieldsToEditorView(p Editable, e *Editor) error {
	defaults := []Field{
		{
			View: Input("Slug", p, map[string]string{
				"label":       "URL Slug",
				"type":        "text",
				"disabled":    "true",
				"placeholder": "Will be set automatically",
			}, nil),
		},
		{
			View: Timestamp("Timestamp", p, map[string]string{
				"type":  "hidden",
				"class": "timestamp __ponzu",
			}),
		},
		{
			View: Timestamp("Updated", p, map[string]string{
				"type":  "hidden",
				"class": "updated __ponzu",
			}),
		},
	}

	for _, f := range defaults {
		err := addFieldToEditorView(e, f)
		if err != nil {
			return err
		}
	}

	return nil
}
