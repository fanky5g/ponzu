package generate

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fanky5g/ponzu/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GenerateTestSuite struct {
	suite.Suite
}

func (suite *GenerateTestSuite) SetupSuite() {
	_, err := os.Stat("./testdata")
	if errors.Is(err, fs.ErrNotExist) {
		err = os.Mkdir("./testdata", os.ModePerm)
	}

	if err != nil {
		suite.T().Fatal(err)
		return
	}
}

func (suite *GenerateTestSuite) TearDownSuite() {
	err := os.RemoveAll("./testdata")
	if err != nil {
		suite.T().Fatal(err)
		return
	}
}

func (suite *GenerateTestSuite) TestWriteTemplateMain() {
	buf := &bytes.Buffer{}
	expectedResult, err := format.Source([]byte(`
package main

import (
	types "github.com/fanky5g/app/entities"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/generator"
	contentGenerator "github.com/fanky5g/ponzu/generator/content"
	modelGenerator "github.com/fanky5g/ponzu/generator/models"
	"log"
)

func main() {
    contentType := generator.Content

	contentTypes := content.Types{
		Content:          types.Content,
		FieldCollections: types.FieldCollections,
		Definitions:      types.Definitions,
	}

    contentGeneratorConfig := generator.Config{
        Target: generator.Target{
            Path: generator.Path{
                Root: "./testdata",
                Base: "entities",
            },
            Package: "entities",
        },
        Type: generator.Content,
    }

    modelGeneratorConfig := generator.Config{
        Target: generator.Target{
            Path: generator.Path{
                Root: "./testdata",
                Base: "models",
            },
            Package: "models",
        },
    }

	contentGeneratorInstance, err := contentGenerator.New(contentGeneratorConfig, contentTypes)
	if err != nil {
		log.Fatal(err)
	}

	generator.Register(contentGeneratorInstance)

	modelGeneratorInstance, err := modelGenerator.New(modelGeneratorConfig, contentGeneratorConfig)
	if err != nil {
        log.Fatal(err)
    }

	generator.Register(modelGeneratorInstance)

	if err = generator.Run(contentType, []string{
		"image_and_text_block:Image and Text Block",
		"image:string image_position:string",
		"content:string:richtext",
		"link:@link",
	}); err != nil {
	log.Panic(err)
}
}
`))

	if err != nil {
		suite.T().Fatal(err)
		return
	}

	if assert.NoError(suite.T(), writeTemplate("main.go.tmpl", map[string]interface{}{
		"ModulePath": "github.com/fanky5g/app",
		"Arguments": []string{
			"image_and_text_block:Image and Text Block",
			"image:string image_position:string",
			"content:string:richtext",
			"link:@link",
		},
		"ContentType": generator.Content,
		"Config": map[string]generator.Config{
			"Content": {
				Target: generator.Target{
					Path: generator.Path{
						Root: "./testdata",
						Base: "entities",
					},
					Package: "entities",
				},
			},
			"Models": {
				Target: generator.Target{
					Path: generator.Path{
						Root: "./testdata",
						Base: "models",
					},
					Package: "models",
				},
			},
		},
	}, buf)) {
		fmt.Println(string(buf.Bytes()))
		assert.Equal(suite.T(), buf.Bytes(), expectedResult)
	}
}

func (suite *GenerateTestSuite) TestWriteTemplateGoMod() {
	buf := &bytes.Buffer{}
	if assert.NoError(suite.T(), writeTemplate("go.mod.tmpl", map[string]interface{}{
		"GoVersion":    "1.16.0",
		"PonzuVersion": "v0.5.1",
		"ModulePath":   "github.com/fanky5g/app",
		"WorkingDir":   "./testdata",
	}, buf)) {
		assert.Equal(suite.T(), buf.String(), `module github.com/fanky5g/app

go 1.16.0

require (
	github.com/fanky5g/ponzu v0.5.1
)
`)
	}
}

func (suite *GenerateTestSuite) writeFile(name string) {
	f, err := os.Create(name)
	if err != nil {
		suite.T().Fatal(err)
		return
	}

	_, err = io.Copy(f, strings.NewReader("Foo Bar"))
	if err != nil {
		suite.T().Fatal(err)
		return
	}

	err = f.Close()
	if err != nil {
		suite.T().Fatal(err)
		return
	}
}

func (suite *GenerateTestSuite) TestCopyFiles() {
	src, err := os.MkdirTemp("./testdata", "src")
	if err != nil {
		suite.T().Fatal(err)
		return
	}

	target, err := os.MkdirTemp("./testdata", "target")
	if err != nil {
		suite.T().Fatal(err)
		return
	}

	suite.writeFile(filepath.Join(src, "dummy.txt"))
	if assert.NoError(suite.T(), copyFiles(src, target)) {
		_, err = os.Stat(filepath.Join(target, "dummy.txt"))
		assert.Equal(suite.T(), nil, err)
	}
}

func TestGenerator(t *testing.T) {
	suite.Run(t, new(GenerateTestSuite))
}
