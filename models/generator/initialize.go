package generator

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/generator"
	"html/template"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func (m *modelGenerator) Initialize(definition *generator.TypeDefinition, writer generator.Writer) error {
	filePath := m.getFileName(definition)

	targetDir := path.Dir(filePath)
	if _, err := os.Stat(targetDir); errors.Is(err, fs.ErrNotExist) {
		var tmpl *template.Template
		tmpl, err = template.ParseFiles(filepath.Join(m.templateDir, "model-root.tmpl"))
		if err != nil {
			return fmt.Errorf("failed to parse template: %s", err.Error())
		}

		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, struct {
			Package string
		}{Package: m.config.Target.Package})
		if err != nil {
			return fmt.Errorf("failed to execute template: %s", err.Error())
		}

		return writer.Write(filepath.Join(targetDir, "ponzu_types.go"), buf.Bytes())
	}

	return nil
}
