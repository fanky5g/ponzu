package editor

import (
	"bytes"
	"github.com/fanky5g/ponzu/config"
	"log"
)

// Form takes editable entities and any number of Field funcs to describe the edit
// page for any entities struct added by a user
func Form(post Editable, paths config.Paths, fields ...Field) ([]byte, error) {
	editor := &Editor{}

	editor.ViewBuf = &bytes.Buffer{}
	_, err := editor.ViewBuf.WriteString(`
		<div class="row">
			<div class="col s12">
			  <ul class="tabs">
				<li class="tab col s3"><a class="active" href="#entities"><i class="material-icons">edit</i>Edit</a></li>
				<li class="tab col s3"><a href="#properties"><i class="material-icons">tune</i>Properties</a></li>
			  </ul>
			</div>
			<div id="entities" class="col s12 editor-entities">
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

	// entities items with Item embedded have some default fields we need to render
	_, err = editor.ViewBuf.WriteString(`<div id="properties" class="col s12 editor-metadata">`)
	if err != nil {
		log.Println("Error writing HTML string to editor Form buffer")
		return nil, err
	}

	contentMetadata := makeHtml("editor_content_metadata")

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
	scriptTmpl := makeScript("editor")
	if err = scriptTmpl.Execute(script, paths); err != nil {
		panic(err)
	}

	editorControls := makeHtml("editor_controls")
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
