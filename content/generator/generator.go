package generator

import (
	content "github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/generator/types"
	"path/filepath"
	"runtime"
)

type Path struct {
	Root string
	Base string
}

type Target struct {
	Path    Path
	Package string
}

type Config struct {
	Types  content.Types
	Target Target
}

type ContentGenerator interface {
	Generate(contentType content.Type, typeDefinition *types.TypeDefinition) error
	ValidateField(field *types.Field) error
}

type generator struct {
	templateDir string
	target      Target
	types       content.Types
}

func setupGenerator(conf Config) (*generator, error) {
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "..")
	templateDir := filepath.Join(rootPath, "generator", "templates")

	return &generator{
		templateDir: templateDir,
		target:      conf.Target,
		types:       conf.Types,
	}, nil
}

func New(conf Config) (ContentGenerator, error) {
	return setupGenerator(conf)
}
