package editor

import (
	"bytes"
	"strings"
)

// Checkbox returns the []byte of a set of <input type="checkbox"> HTML elements
// wrapped in a <div> with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Checkbox(fieldName string, p interface{}, attrs, options map[string]string) []byte {
	if _, ok := attrs["class"]; ok {
		attrs["class"] += "input-field col s12"
	} else {
		attrs["class"] = "input-field col s12"
	}

	div := NewElement("div", attrs["label"], fieldName, p, attrs, nil)

	var opts []*Element

	// get the pre-checked options if this is already an existing post
	checkedVals := ValueFromStructField(fieldName, p, nil).(string)
	checked := strings.Split(checkedVals, "__ponzu")

	i := 0
	for k, v := range options {
		inputAttrs := map[string]string{
			"type":  "checkbox",
			"value": k,
			"id":    strings.Join(strings.Split(v, " "), "-"),
		}

		// check if k is in the pre-checked values and set to checked
		for _, x := range checked {
			if k == x {
				inputAttrs["checked"] = "checked"
			}
		}

		// create a *element manually using the modified TagNameFromStructFieldMulti
		// func since this is for a multi-value name
		input := &Element{
			TagName: "input",
			Attrs:   inputAttrs,
			Name:    TagNameFromStructFieldMulti(fieldName, i, p),
			Label:   v,
			Data:    "",
			ViewBuf: &bytes.Buffer{},
		}

		opts = append(opts, input)
		i++
	}

	return DOMElementWithChildrenCheckbox(div, opts)
}
