package generator

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/generator"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func (gt *contentGenerator) getFileName(definition *generator.TypeDefinition) string {
	fileName := strings.ToLower(definition.Name) + ".go"
	return filepath.Join(gt.config.Target.Path.Root, gt.config.Target.Path.Base, fileName)
}

func (gt *contentGenerator) Generate(definition *generator.TypeDefinition, writer generator.Writer) error {
	filePath := gt.getFileName(definition)
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		localFile := filepath.Join(gt.config.Target.Path.Root, gt.config.Target.Path.Base, filePath)
		return fmt.Errorf("please remove '%s' before executing this command", localFile)
	}

	templateName, err := gt.getTemplateName(definition)
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(templateName)
	if err != nil {
		return fmt.Errorf("failed to parse template: %s", err.Error())
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, newViewScope(definition, gt.contentTypes, gt.config.Target, gt.templateDir))

	if err != nil {
		return fmt.Errorf("failed to execute template: %s", err.Error())
	}

	return writer.Write(filePath, buf.Bytes())
}

func (gt *contentGenerator) getTemplateName(definition *generator.TypeDefinition) (string, error) {
	switch definition.Type {
	case generator.Content:
		return filepath.Join(gt.templateDir, "gen-content.tmpl"), nil
	case generator.Plain:
		return filepath.Join(gt.templateDir, "gen-type.tmpl"), nil
	case generator.FieldCollection:
		return filepath.Join(gt.templateDir, "gen-field-collection.tmpl"), nil
	default:
		return "", errors.New("unsupported content type")
	}
}
