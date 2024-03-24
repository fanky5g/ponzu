package generate

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	contentGenerator "github.com/fanky5g/ponzu/content/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go/format"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type GenerateTestSuite struct {
	suite.Suite
	g Generator
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

	suite.g, err = NewGenerator(GeneratorArgs{
		ContentType: content.TypePlain,
		Arguments: []string{
			"image_and_text_block:Image and Text Block",
			"image:string image_position:string",
			"content:string:richtext",
			"link:@link",
		},
		Target: contentGenerator.Target{
			Path: contentGenerator.Path{
				Root: "./testdata",
				Base: "entities",
			},
			Package: "entities",
		},
	})

	if err != nil {
		suite.T().Fatal(err)
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
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/generator"
	"github.com/fanky5g/ponzu/content/generator/parser"
	targetPackage "content-generator/entities"
	"log"
)

func main() {
	types := content.Types{
		Content:          targetPackage.Content,
		FieldCollections: targetPackage.FieldCollections,
		Definitions:      targetPackage.Definitions,
	}

	p, err := parser.New(types)

	// parse type info from args
	typeDefinition, err := p.ParseTypeDefinition(content.TypePlain, []string{
		"image_and_text_block:Image and Text Block",
		"image:string image_position:string",
		"content:string:richtext",
		"link:@link",
	})

	if err != nil {
		log.Panicf("failed to parse type args: %s\n", err.Error())
	}

	domainContentGenerator, err := generator.New(generator.Config{
		Types: types,
		Target: generator.Target{
			Path: generator.Path{
				Root: "./testdata",
				Base: "entities",
			},
			Package: "entities",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	contentGenerators := []generator.ContentGenerator{domainContentGenerator}
	for _, contentGenerator := range contentGenerators {
		for _, field := range typeDefinition.Fields {
			if err = contentGenerator.ValidateField(&field); err != nil {
				log.Panic(err)
			}
		}

		if err = contentGenerator.Generate(content.TypePlain, typeDefinition); err != nil {
			log.Panic(err)
		}
	}
}
`))

	if err != nil {
		suite.T().Fatal(err)
		return
	}

	if assert.NoError(suite.T(), suite.g.WriteTemplate("main.go.tmpl", buf)) {
		assert.Equal(suite.T(), buf.Bytes(), expectedResult)
	}
}

func (suite *GenerateTestSuite) TestWriteTemplateGoMod() {
	buf := &bytes.Buffer{}
	if assert.NoError(suite.T(), suite.g.WriteTemplate("go.mod.tmpl", buf)) {
		assert.Equal(suite.T(), buf.String(), fmt.Sprintf(`module content-generator

go %s`, strings.TrimPrefix(runtime.Version(), "go")))
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
	if assert.NoError(suite.T(), suite.g.CopyFiles(src, target)) {
		_, err = os.Stat(filepath.Join(target, "dummy.txt"))
		assert.Equal(suite.T(), nil, err)
	}
}

func TestGenerator(t *testing.T) {
	suite.Run(t, new(GenerateTestSuite))
}
