package editor

import (
	"bytes"
	"html"
)

// Richtext returns the []byte of a rich text editor (provided by http://summernote.org/) with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Richtext(fieldName string, p interface{}, attrs map[string]string, args *FieldArgs) []byte {
	// create wrapper for richtext editor, which isolates the editor's css
	iso := []byte(`<div class="iso-texteditor control-block"><label>` + attrs["label"] + `</label>`)
	isoClose := []byte(`</div>`)

	name := TagNameFromStructField(fieldName, p, args)
	if _, ok := attrs["class"]; ok {
		attrs["class"] += "richtext " + fieldName
	} else {
		attrs["class"] = "richtext " + fieldName
	}

	if _, ok := attrs["id"]; ok {
		attrs["id"] += "richtext-" + name
	} else {
		attrs["id"] = "richtext-" + name
	}

	// create the target element for the editor to attach itself
	div := &Element{
		TagName: "div",
		Attrs:   attrs,
		Name:    "",
		Label:   "",
		Data:    "",
		ViewBuf: &bytes.Buffer{},
	}

	// create a hidden input to store the value from the struct
	val := ValueFromStructField(fieldName, p, args).(string)
	input := `<input type="hidden" name="` + name + `" class="richtext-value ` + fieldName + `" value="` + html.EscapeString(val) + `"/>`

	// build the dom tree for the entire richtext component
	iso = append(iso, DOMElement(div)...)
	iso = append(iso, []byte(input)...)
	iso = append(iso, isoClose...)

	script := &bytes.Buffer{}
	scriptTmpl := makeScript("richtext")

	if err := scriptTmpl.Execute(script, struct {
		InputName string
		Attrs     map[string]string
	}{
		InputName: name,
		Attrs:     attrs,
	}); err != nil {
		panic(err)
	}

	return append(iso, script.Bytes()...)
}
