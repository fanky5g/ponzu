package editor

import "bytes"

// Timestamp returns the []byte of an <input> HTML element with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Timestamp(fieldName string, p interface{}, attrs map[string]string) []byte {
	var data string
	val := ValueFromStructField(fieldName, p, nil).(string)
	if val == "0" {
		data = ""
	} else {
		data = val
	}

	e := &Element{
		TagName: "input",
		Attrs:   attrs,
		Name:    TagNameFromStructField(fieldName, p, nil),
		Label:   attrs["label"],
		Data:    data,
		ViewBuf: &bytes.Buffer{},
	}

	return DOMInputSelfClose(e)
}
