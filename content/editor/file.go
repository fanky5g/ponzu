package editor

import (
	"bytes"
	"github.com/fanky5g/ponzu/internal/templates"
	"net/url"
)

const PonzuFileStorageRoute = "/api/uploads"

// File returns the []byte of a <input type="file"> HTML element with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func File(publicPath, fieldName string, p interface{}, attrs map[string]string) []byte {
	name := TagNameFromStructField(fieldName, p, nil)
	path := ValueFromStructField(fieldName, p, nil).(string)

	externalPath, err := url.JoinPath(PonzuFileStorageRoute, path)
	if err != nil {
		panic(err)
	}

	w := &bytes.Buffer{}
	file := struct {
		Name       string
		Path       string
		Attributes map[string]string
		PublicPath string
		URL        string
	}{
		Name:       name,
		Path:       path,
		Attributes: attrs,
		PublicPath: publicPath,
		URL:        externalPath,
	}

	if err = templates.ExecuteTemplate(w, "views/input/file.gohtml", file); err != nil {
		panic(err)
	}

	return w.Bytes()
}
