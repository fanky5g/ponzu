package generator

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/generator/types"
	log "github.com/sirupsen/logrus"
	"go/format"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

func (gt *generator) getTemplateName(contentType content.Type) (string, error) {
	switch contentType {
	case content.TypeContent:
		return filepath.Join(gt.templateDir, "gen-content.tmpl"), nil
	case content.TypePlain:
		return filepath.Join(gt.templateDir, "gen-type.tmpl"), nil
	case content.TypeFieldCollection:
		return filepath.Join(gt.templateDir, "gen-field-collection.tmpl"), nil
	}

	return "", errors.New("unsupported content type")
}

func (gt *generator) Generate(contentType content.Type, definition *types.TypeDefinition) error {
	fileName := strings.ToLower(definition.Name) + ".go"
	filePath := filepath.Join(gt.target.Path.Root, gt.target.Path.Base, fileName)

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		localFile := filepath.Join(gt.target.Path.Root, gt.target.Path.Base, fileName)
		return fmt.Errorf("please remove '%s' before executing this command", localFile)
	}

	targetDir := path.Dir(filePath)
	if _, err := os.Stat(targetDir); errors.Is(err, fs.ErrNotExist) {
		if err = gt.makeRootContentDir(targetDir); err != nil {
			return err
		}
	}

	buf, err := gt.generate(contentType, definition)
	if err != nil {
		return err
	}

	return gt.writeFile(filePath, buf)
}

func (gt *generator) writeFile(filePath string, buf []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}(file)

	_, err = file.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write generated file buffer: %s", err.Error())
	}

	return nil
}

func (gt *generator) makeRootContentDir(target string) error {
	if err := os.MkdirAll(target, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create content director: %v", err)
	}

	tmpl, err := template.ParseFiles(filepath.Join(gt.templateDir, "content-root.tmpl"))
	if err != nil {
		return fmt.Errorf("failed to parse template: %s", err.Error())
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, struct {
		Package string
	}{Package: gt.target.Package})
	if err != nil {
		return fmt.Errorf("failed to execute template: %s", err.Error())
	}

	return gt.writeFile(filepath.Join(target, "ponzu_types.go"), buf.Bytes())
}

func (gt *generator) generate(contentType content.Type, definition *types.TypeDefinition) ([]byte, error) {
	for i := range definition.Fields {
		if err := gt.setFieldView(definition, i); err != nil {
			return nil, err
		}
	}

	tmplPath, err := gt.getTemplateName(contentType)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %s", err.Error())
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, struct {
		Definition *types.TypeDefinition
		Target     Target
	}{Definition: definition, Target: gt.target})
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %s", err.Error())
	}

	fmtBuf, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format template: %s", err.Error())
	}

	return fmtBuf, nil
}
