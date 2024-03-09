package editor

import "bytes"

// Select returns the []byte of a <select> HTML element plus internal <options> with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Select(fieldName string, p interface{}, attrs, options map[string]string) []byte {
	// options are the value attr and the display value, i.e.
	// <option value="{map key}">{map value}</option>

	// find the field value in p to determine if an option is pre-selected
	fieldVal := ValueFromStructField(fieldName, p, nil)

	if _, ok := attrs["class"]; ok {
		attrs["class"] += " browser-default"
	} else {
		attrs["class"] = "browser-default"
	}

	sel := NewElement("select", attrs["label"], fieldName, p, attrs, nil)
	var opts []*Element

	// provide a call to action for the select element
	cta := &Element{
		TagName: "option",
		Attrs:   map[string]string{"disabled": "true", "selected": "true"},
		Data:    "Select an option...",
		ViewBuf: &bytes.Buffer{},
	}

	// provide a selection reset (will store empty string in db)
	reset := &Element{
		TagName: "option",
		Attrs:   map[string]string{"value": ""},
		Data:    "None",
		ViewBuf: &bytes.Buffer{},
	}

	opts = append(opts, cta, reset)

	for k, v := range options {
		optAttrs := map[string]string{"value": k}
		if k == fieldVal {
			optAttrs["selected"] = "true"
		}
		opt := &Element{
			TagName: "option",
			Attrs:   optAttrs,
			Data:    v,
			ViewBuf: &bytes.Buffer{},
		}

		opts = append(opts, opt)
	}

	return DOMElementWithChildrenSelect(sel, opts)
}
