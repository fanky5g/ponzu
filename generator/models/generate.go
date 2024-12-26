package models

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/generator"
	"github.com/fanky5g/ponzu/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"text/template"
)

var generateFuncMap = template.FuncMap{
	"TableName": func(name string) string {
		s, err := util.Slugify(name)
		if err != nil {
			log.WithField("Error", err).Panic("Failed to build table name")
		}

		return strings.Replace(s, "-", "_", -1)
	},
}

func (m *modelGenerator) Generate(definition *generator.TypeDefinition, writer generator.Writer) error {
	if definition.Type != generator.Content {
		return nil
	}

	fileName := m.getFileName(definition)

	var tmpl *template.Template
	tmpl, err := template.New("gen-model.tmpl").Funcs(generateFuncMap).ParseFiles(filepath.Join(m.templateDir, "gen-model.tmpl"))
	if err != nil {
		return fmt.Errorf("failed to parse template: %s", err.Error())
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(
		buf,
		struct {
			ContentPackage string
			ContentRoot    string
			ContentBase    string
			Definition     *generator.TypeDefinition
			Package        string
			ModulePath     string
		}{
			ContentPackage: m.contentGeneratorConfig.Target.Package,
			ContentBase:    m.contentGeneratorConfig.Target.Path.Base,
			Definition:     definition,
			Package:        m.config.Target.Package,
			ModulePath:     m.modulePath,
		},
	)

	if err != nil {
		return errors.Wrap(err, "Failed to execute template")
	}

	return writer.Write(fileName, buf.Bytes())
}
