package editor

import (
	"bytes"
	"html"
	"log"
	"strings"
)

// Element is a basic struct for representing DOM elements
type Element struct {
	TagName string
	Attrs   map[string]string
	Name    string
	Label   string
	Data    string
	ViewBuf *bytes.Buffer
}

// NewElement returns an Element with Name and Data already processed from the
// fieldName and entities interface provided
func NewElement(tagName, label, fieldName string, p interface{}, attrs map[string]string, args *FieldArgs) *Element {
	if attrs == nil {
		attrs = make(map[string]string)
	}

	// define mdc element attributes
	if tagName == "input" {
		addClassName(attrs, "mdc-text-field__input")
	}

	return &Element{
		TagName: tagName,
		Attrs:   attrs,
		Name:    TagNameFromStructField(fieldName, p, args),
		Label:   label,
		Data:    ValueFromStructField(fieldName, p, args).(string),
		ViewBuf: &bytes.Buffer{},
	}
}

func DOMInputSelfClose(e *Element) []byte {
	isValidSelfClosingTag := true
	if len(e.Attrs) > 0 {
		if elementType, ok := e.Attrs["type"]; ok {
			if elementType == "hidden" {
				isValidSelfClosingTag = false
			}
		}
	}

	var err error
	if isValidSelfClosingTag {
		_, err = e.ViewBuf.WriteString(`<fieldset class="control-block">`)
		if err != nil {
			log.Println("Error writing HTML string to buffer: DOMElementSelfClose")
			return nil
		}

		_, err = e.ViewBuf.WriteString(`
            <label class="mdc-text-field mdc-text-field--filled">
              <span class="mdc-text-field__ripple"></span>
              <span class="mdc-floating-label" id="my-label-id">` + e.Label + `</span>
    `)

		if err != nil {
			log.Println("Error writing HTML string to buffer: DOMElementSelfClose")
			return nil
		}
	}

	_, err = e.ViewBuf.WriteString(`<` + e.TagName + ` value="`)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElementSelfClose")
		return nil
	}

	_, err = e.ViewBuf.WriteString(html.EscapeString(e.Data) + `" `)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElementSelfClose")
		return nil
	}

	for attr, value := range e.Attrs {
		_, err := e.ViewBuf.WriteString(attr + `="` + value + `" `)
		if err != nil {
			log.Println("Error writing HTML string to buffer: DOMElementSelfClose")
			return nil
		}
	}
	_, err = e.ViewBuf.WriteString(` name="` + e.Name + `" />`)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElementSelfClose")
		return nil
	}

	if isValidSelfClosingTag {
		_, err = e.ViewBuf.WriteString(`<span class="mdc-line-ripple"></span></label></fieldset>`)
		if err != nil {
			log.Println("Error writing HTML string to buffer: DOMElementSelfClose")
			return nil
		}
	}

	return e.ViewBuf.Bytes()
}

// DOMElement creates a DOM element
func DOMElement(e *Element) []byte {
	_, err := e.ViewBuf.WriteString(`<div class="input-field col s12">`)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElement")
		return nil
	}

	if e.Label != "" {
		_, err = e.ViewBuf.WriteString(
			`<label class="active" for="` +
				strings.Join(strings.Split(e.Label, " "), "-") + `">` + e.Label +
				`</label>`)
		if err != nil {
			log.Println("Error writing HTML string to buffer: DOMElement")
			return nil
		}
	}

	_, err = e.ViewBuf.WriteString(`<` + e.TagName + ` `)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElement")
		return nil
	}

	for attr, value := range e.Attrs {
		_, err = e.ViewBuf.WriteString(attr + `="` + string(value) + `" `)
		if err != nil {
			log.Println("Error writing HTML string to buffer: DOMElement")
			return nil
		}
	}
	_, err = e.ViewBuf.WriteString(` name="` + e.Name + `" >`)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElement")
		return nil
	}

	_, err = e.ViewBuf.WriteString(html.EscapeString(e.Data))
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElement")
		return nil
	}

	_, err = e.ViewBuf.WriteString(`</` + e.TagName + `>`)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElement")
		return nil
	}

	_, err = e.ViewBuf.WriteString(`</div>`)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElement")
		return nil
	}

	return e.ViewBuf.Bytes()
}

func DOMElementWithChildrenSelect(e *Element, children []*Element) []byte {
	_, err := e.ViewBuf.WriteString(`<div class="control-block">`)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElementWithChildrenSelect")
		return nil
	}

	_, err = e.ViewBuf.WriteString(`<` + e.TagName + ` `)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElementWithChildrenSelect")
		return nil
	}

	for attr, value := range e.Attrs {
		_, err = e.ViewBuf.WriteString(attr + `="` + value + `" `)
		if err != nil {
			log.Println("Error writing HTML string to buffer: DOMElementWithChildrenSelect")
			return nil
		}
	}
	_, err = e.ViewBuf.WriteString(` name="` + e.Name + `" >`)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElementWithChildrenSelect")
		return nil
	}

	// loop over children and create DOMElement for each child
	for _, child := range children {
		_, err = e.ViewBuf.Write(DOMElement(child))
		if err != nil {
			log.Println("Error writing HTML DOMElement to buffer: DOMElementWithChildrenSelect")
			return nil
		}
	}

	_, err = e.ViewBuf.WriteString(`</` + e.TagName + `>`)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElementWithChildrenSelect")
		return nil
	}

	if e.Label != "" {
		_, err = e.ViewBuf.WriteString(`<label class="active">` + e.Label + `</label>`)
		if err != nil {
			log.Println("Error writing HTML string to buffer: DOMElementWithChildrenSelect")
			return nil
		}
	}

	_, err = e.ViewBuf.WriteString(`</div>`)
	if err != nil {
		log.Println("Error writing HTML string to buffer: DOMElementWithChildrenSelect")
		return nil
	}

	return e.ViewBuf.Bytes()
}

func addClassName(attrs map[string]string, className string) {
	if attrs == nil || strings.TrimSpace(className) == "" {
		return
	}

	if _, ok := attrs["class"]; !ok {
		attrs["class"] = ""
	}

	classNames := strings.FieldsFunc(attrs["class"], func(r rune) bool {
		return r == ' '
	})

	hasClassName := false
	for _, class := range classNames {
		if strings.TrimSpace(class) == className {
			hasClassName = true
			break
		}
	}

	if !hasClassName {
		classNames = append(classNames, className)
	}

	attrs["class"] = strings.Join(classNames, " ")
}
