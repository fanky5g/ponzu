package editor

import (
	"bytes"
	"log"
	"strings"
)

// SelectRepeater returns the []byte of a <select> HTML element plus internal <options> with a label.
// It also includes repeat controllers (+ / -) so the element can be
// dynamically multiplied or reduced.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func SelectRepeater(fieldName string, p interface{}, attrs, options map[string]string) []byte {
	// options are the value attr and the display value, i.e.
	// <option value="{map key}">{map value}</option>
	scope := TagNameFromStructField(fieldName, p, nil)
	html := bytes.Buffer{}
	_, err := html.WriteString(`<span class="__ponzu-repeat ` + scope + `">`)
	if err != nil {
		log.Println("Error writing HTML string to SelectRepeater buffer")
		return nil
	}

	// find the field values in p to determine if an option is pre-selected
	fieldVals := ValueFromStructField(fieldName, p, nil).(string)
	vals := strings.Split(fieldVals, "__ponzu")

	if _, ok := attrs["class"]; ok {
		attrs["class"] += " browser-default"
	} else {
		attrs["class"] = "browser-default"
	}

	// loop through vals and create selects and options for each, adding to views
	if len(vals) > 0 {
		for i, val := range vals {
			sel := &Element{
				TagName: "select",
				Attrs:   attrs,
				Name:    TagNameFromStructFieldMulti(fieldName, i, p),
				ViewBuf: &bytes.Buffer{},
			}

			// only add the label to the first select in repeated list
			if i == 0 {
				sel.Label = attrs["label"]
			}

			// create options for select element
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
				if k == val {
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

			_, err := html.Write(DOMElementWithChildrenSelect(sel, opts))
			if err != nil {
				log.Println("Error writing DOMElementWithChildrenSelect to SelectRepeater buffer")
				return nil
			}
		}
	}

	_, err = html.WriteString(`</span>`)
	if err != nil {
		log.Println("Error writing HTML string to SelectRepeater buffer")
		return nil
	}

	return append(html.Bytes(), RepeatController(fieldName, p, "select", ".input-field")...)
}
