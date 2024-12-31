package content

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/generator"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

func (gt *contentGenerator) Initialize(definition *generator.TypeDefinition, writer generator.Writer) error {
	filePath := gt.getFileName(definition)

	targetDir := path.Dir(filePath)
	if _, err := os.Stat(targetDir); errors.Is(err, fs.ErrNotExist) {
		if err = os.MkdirAll(targetDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create content director: %v", err)
		}

		var tmpl *template.Template
		tmpl, err = template.ParseFiles(filepath.Join(gt.templateDir, "content-root.tmpl"))
		if err != nil {
			return fmt.Errorf("failed to parse template: %s", err.Error())
		}

		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, struct {
			Package string
		}{Package: gt.config.Target.Package})
		if err != nil {
			return fmt.Errorf("failed to execute template: %s", err.Error())
		}

		return writer.Write(filepath.Join(targetDir, "ponzu_types.go"), buf.Bytes())
	}

	return nil
}
