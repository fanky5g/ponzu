package generator

import (
	"fmt"
	"github.com/fanky5g/ponzu/generator"
	"github.com/stretchr/testify/suite"
	"go/format"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type testWriter struct{ buf []byte }

func (writer *testWriter) Write(filePath string, buf []byte) error {
	writer.buf = buf
	var err error
	if strings.HasSuffix(strings.TrimSuffix(filePath, ".tmpl"), ".go") {
		writer.buf, err = format.Source(buf)
		if err != nil {
			return fmt.Errorf("failed to format template: %s", err.Error())
		}
	}

	return nil
}

type GeneratorTestSuite struct {
	suite.Suite
	g generator.Generator
}

func (s *GeneratorTestSuite) SetupSuite() {
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../..")

	m := &modelGenerator{
		modulePath:  "github.com/fanky5g/testapp",
		templateDir: filepath.Join(rootPath, "models", "generator", "templates"),
		config: generator.Config{
			Target: generator.Target{
				Path: generator.Path{
					Base: "models",
				},
				Package: "models",
			},
		},
		contentGeneratorConfig: generator.Config{
			Target: generator.Target{
				Path: generator.Path{
					Base: "entities",
				},
				Package: "entities",
			},
		},
	}

	s.g = m
}

func TestGenerator(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}
