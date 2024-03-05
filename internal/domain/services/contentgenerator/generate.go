package contentgenerator

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/domain/enum"
	"github.com/fanky5g/ponzu/internal/domain/interfaces"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

type generator struct {
	templateDir string
	contentDir  string
}

func (gt *generator) getTemplateName(contentType enum.ContentType) (string, error) {
	switch contentType {
	case enum.TypeContent:
		return filepath.Join(gt.templateDir, "gen-content.tmpl"), nil
	case enum.TypePlain:
		return filepath.Join(gt.templateDir, "gen-type.tmpl"), nil
	case enum.TypeFieldCollection:
		return filepath.Join(gt.templateDir, "gen-field-collection.tmpl"), nil
	}

	return "", errors.New("unsupported content type")
}

func (gt *generator) Generate(contentType enum.ContentType, definition *item.TypeDefinition) error {
	fileName := strings.ToLower(definition.Name) + ".go"
	filePath := filepath.Join(gt.contentDir, fileName)

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		localFile := filepath.Join(gt.contentDir, fileName)
		return fmt.Errorf("please remove '%s' before executing this command", localFile)
	}

	buf, err := gt.generate(contentType, definition)
	if err != nil {
		return err
	}

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

func (gt *generator) generate(contentType enum.ContentType, definition *item.TypeDefinition) ([]byte, error) {
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
	err = tmpl.Execute(buf, definition)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %s", err.Error())
	}

	fmtBuf, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format template: %s", err.Error())
	}

	return fmtBuf, nil
}

func setupGenerator() (*generator, error) {
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../../../..")
	templateDir := filepath.Join(rootPath, "internal", "domain", "services", "contentgenerator", "templates")
	contentDir := filepath.Join(rootPath, "internal", "domain", "entities", "content")

	return &generator{
		templateDir: templateDir,
		contentDir:  contentDir,
	}, nil
}

func New() (interfaces.ContentGenerator, error) {
	return setupGenerator()
}
