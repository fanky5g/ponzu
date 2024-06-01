package generator

import (
	"github.com/fanky5g/ponzu/generator"
	"github.com/fanky5g/ponzu/util"
	"path/filepath"
	"runtime"
	"strings"
)

type modelGenerator struct {
	modulePath             string
	templateDir            string
	config                 generator.Config
	contentGeneratorConfig generator.Config
}

func (m *modelGenerator) getFileName(definition *generator.TypeDefinition) string {
	fileName := strings.ToLower(definition.Name) + ".go"
	return filepath.Join(m.config.Target.Path.Root, m.config.Target.Path.Base, fileName)
}

func New(
	config generator.Config,
	contentGeneratorConfig generator.Config) (generator.Generator, error) {
	modulePath, err := util.GetModulePath()
	if err != nil {
		return nil, err
	}

	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../..")

	return &modelGenerator{
		modulePath:             modulePath,
		templateDir:            filepath.Join(rootPath, "models", "generator", "templates"),
		config:                 config,
		contentGeneratorConfig: contentGeneratorConfig,
	}, nil
}
