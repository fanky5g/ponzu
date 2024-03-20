package editor

import (
	"bytes"
	"log"
	"strings"
)

func InputRepeater(fieldName string, p interface{}, attrs map[string]string) []byte {
	// find the field values in p to determine pre-filled inputs
	fieldVals := ValueFromStructField(fieldName, p, nil).(string)
	vals := strings.Split(fieldVals, "__ponzu")

	scope := TagNameFromStructField(fieldName, p, nil)
	html := bytes.Buffer{}

	_, err := html.WriteString(`<span class="__ponzu-repeat ` + scope + `">`)
	if err != nil {
		log.Println("Error writing HTML string to InputRepeater buffer")
		return nil
	}

	for i, val := range vals {
		el := &Element{
			TagName: "input",
			Attrs:   attrs,
			Name:    TagNameFromStructFieldMulti(fieldName, i, p),
			Data:    val,
			ViewBuf: &bytes.Buffer{},
		}

		// only add the label to the first input in repeated list
		if i == 0 {
			el.Label = attrs["label"]
		}

		_, err := html.Write(DOMElementSelfClose(el))
		if err != nil {
			log.Println("Error writing DOMElementSelfClose to InputRepeater buffer")
			return nil
		}
	}
	_, err = html.WriteString(`</span>`)
	if err != nil {
		log.Println("Error writing HTML string to InputRepeater buffer")
		return nil
	}

	return append(html.Bytes(), RepeatController(fieldName, p, "input", ".input-field")...)
}
