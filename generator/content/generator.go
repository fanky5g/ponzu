package content

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/generator"
	"path/filepath"
	"runtime"
)

type contentGenerator struct {
	templateDir  string
	contentTypes content.Types
	config       generator.Config
}

func setupGenerator(config generator.Config, contentTypes content.Types) (*contentGenerator, error) {
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../..")
	templateDir := filepath.Join(rootPath, "generator", "content", "templates")

	return &contentGenerator{
		templateDir:  templateDir,
		contentTypes: contentTypes,
		config:       config,
	}, nil
}

func New(config generator.Config, contentTypes content.Types) (generator.Generator, error) {
	return setupGenerator(config, contentTypes)
}
