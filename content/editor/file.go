package editor

import (
	"bytes"
	"strings"

	"github.com/fanky5g/ponzu/internal/templates"
)

// File returns the []byte of a <input type="file"> HTML element with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func File(fieldName string, p interface{}, attrs map[string]string) []byte {
	name := TagNameFromStructField(fieldName, p, nil)
	value := ValueFromStructField(fieldName, p, nil).(string)

	publicPath := ""
	if publicPathValue, ok := attrs["PublicPath"]; ok {
		publicPath = publicPathValue
	}

	w := &bytes.Buffer{}
	file := struct {
		Name         string
		Path         string
		RelativePath string
		Attributes   map[string]string
		PublicPath   string
	}{
		Name:         name,
		Path:         strings.TrimPrefix(value, publicPath),
		RelativePath: value,
		Attributes:   attrs,
		PublicPath:   publicPath,
	}
	if err := templates.ExecuteTemplate(w, "views/input/file.gohtml", file); err != nil {
		panic(err)
	}

	return w.Bytes()
}
