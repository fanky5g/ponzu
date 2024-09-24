package editor

import (
	"bytes"
	"log"
	"strings"
)

func DOMCheckboxSelfClose(e *Element) []byte {
	_, err := e.ViewBuf.WriteString(`
    <fieldset class="control-block">
        <div class="mdc-form-field">
          <div class="mdc-checkbox">
    `)

	// write checkbox control
	addClassName(e.Attrs, "mdc-checkbox__native-control")
	_, err = e.ViewBuf.WriteString(`<` + e.TagName + ` `)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElementCheckbox")
		return nil
	}

	for attr, value := range e.Attrs {
		_, err := e.ViewBuf.WriteString(attr + `="` + value + `" `)
		if err != nil {
			log.Println("Error writing HTML string to buffer: DOMElementCheckbox")
			return nil
		}
	}

	_, err = e.ViewBuf.WriteString(` name="` + e.Name + `" />`)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElementCheckbox")
		return nil
	}
	// end checkbox control

	_, err = e.ViewBuf.WriteString(`
            <div class="mdc-checkbox__background">
              <svg class="mdc-checkbox__checkmark"
                   viewBox="0 0 24 24">
                <path class="mdc-checkbox__checkmark-path"
                      fill="none"
                      d="M1.73,12.91 8.1,19.28 22.79,4.59"/>
              </svg>
              <div class="mdc-checkbox__mixedmark"></div>
            </div>
            <div class="mdc-checkbox__ripple"></div>
          </div>
          <label for="` + e.Attrs["id"] + `">` + e.Label + `</label>
        </div>
    </fieldset>
    `)

	return e.ViewBuf.Bytes()
}

// Checkbox returns the []byte of a set of <input type="checkbox"> HTML elements
// wrapped in a <div> with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Checkbox(fieldName string, p interface{}, attrs map[string]string) []byte {
	value := ValueFromStructField(fieldName, p, nil).(string)
	checked := value == "true"

	inputAttrs := map[string]string{
		"type":  "checkbox",
		"id":    strings.Join(strings.Split(attrs["label"], " "), "-"),
		"value": value,
	}

	if checked {
		inputAttrs["checked"] = "checked"
	}

	// create a *element manually using the modified TagNameFromStructFieldMulti
	// func since this is for a multi-value name
	input := &Element{
		TagName: "input",
		Attrs:   inputAttrs,
		// TODO: make checkbox work without indexing
		Name:    TagNameFromStructFieldMulti(fieldName, 0, p),
		Label:   attrs["label"],
		Data:    "",
		ViewBuf: &bytes.Buffer{},
	}

	return DOMCheckboxSelfClose(input)
}
