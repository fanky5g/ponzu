package generator

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/generator/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go/format"
	"testing"
)

type GenerateTestSuite struct {
	suite.Suite
	gt *generator
}

func (s *GenerateTestSuite) SetupSuite() {
	var err error
	s.gt, err = setupGenerator(Config{
		Types: content.Types{},
		Target: Target{
			Path: Path{
				Root: "",
				Base: "",
			},
			Package: "entities",
		},
	})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *GenerateTestSuite) TestGenerateFieldCollection() {
	t := &types.TypeDefinition{
		Name:    "PageContentBlocks",
		Label:   "Page Content Blocks",
		Initial: "p",
		ContentBlocks: []types.ContentBlock{
			{
				TypeName: "ImageGallery",
				Label:    "Image Gallery",
			},
			{
				TypeName: "ImageAndTextBlock",
				Label:    "Image And Text Block",
			},
			{
				TypeName: "TextBlock",
				Label:    "Text Block",
			},
		},
	}

	expectedBuffer, err := format.Source([]byte(`
	package entities

	import (
		"encoding/json"
		"fmt"
		"github.com/fanky5g/ponzu/content"
	)

	type PageContentBlocks []content.FieldCollection

	func (p *PageContentBlocks) Name() string {
		return "Page Content Blocks"
	}

	func (p *PageContentBlocks) Data() []content.FieldCollection {
		if p == nil {
			return nil
		}
	
		return *p
	}

	func (p *PageContentBlocks) AllowedTypes() map[string]content.Builder {
		return map[string]content.Builder{
			"ImageGallery": func() interface{} {
				return new(ImageGallery)
			},
			"ImageAndTextBlock": func() interface{} {
				return new(ImageAndTextBlock)
			},
			"TextBlock": func() interface{} {
				return new(TextBlock)
			},
		}
	}

	func (p *PageContentBlocks) Add(fieldCollection content.FieldCollection) {
		*p = append(*p, fieldCollection)
	}

	func (p *PageContentBlocks) Set(i int, fieldCollection content.FieldCollection) {
		data := p.Data()
		data[i] = fieldCollection
		*p = data
	}

	func (p *PageContentBlocks) SetData(data []content.FieldCollection) {
		*p = data
	}

	func (p *PageContentBlocks) UnmarshalJSON(b []byte) error {
		if p == nil {
			*p = make([]content.FieldCollection, 0)
		}

		allowedTypes := p.AllowedTypes()

		var value []content.FieldCollection
		if err := json.Unmarshal(b, &value); err != nil {
			return err
		}

		for i, t := range value {
			builder, ok := allowedTypes[t.Type]
			if !ok {
				return fmt.Errorf("type %s not implemented", t.Type)
			}

			entity := builder()
			byteRepresentation, err := json.Marshal(t.Value)
			if err != nil {
				return err
			}

			if err = json.Unmarshal(byteRepresentation, entity); err != nil {
				return err
			}

			value[i].Value = entity
		}

		*p = value
		return nil
	}

	func init() {
		FieldCollections["PageContentBlocks"] = func() interface{} {
			return new(PageContentBlocks)
		}
	}
	`))

	if err != nil {
		s.T().Fatal(err)
	}

	buf, err := s.gt.generate(content.TypeFieldCollection, t)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expectedBuffer, buf)
	}
}

func (s *GenerateTestSuite) TestGenerateType() {
	t := &types.TypeDefinition{
		Name:    "Author",
		Label:   "Author",
		Initial: "a",
		Fields: []types.Field{
			{
				Name:              "Name",
				Label:             "Name",
				Initial:           "a",
				TypeName:          "string",
				JSONName:          "name",
				ViewType:          "input",
				IsReference:       false,
				IsNested:          false,
				ReferenceName:     "",
				ReferenceJSONTags: []string{},
			},
			{
				Name:              "Age",
				Label:             "Age",
				Initial:           "a",
				TypeName:          "int",
				JSONName:          "age",
				ViewType:          "input",
				IsReference:       false,
				IsNested:          false,
				ReferenceName:     "",
				ReferenceJSONTags: []string{},
			},
		},
		HasReferences: false,
	}

	expectedBuffer, err := format.Source([]byte(`
	package entities

import (
	"github.com/fanky5g/ponzu/content"
)

type Author struct {
	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
}

func init() {
	Definitions["Author"] = content.TypeDefinition{
		Name:    "Author",
		Label:   "Author",
		Initial: "a",
		Fields: []content.Field{
			{
				Name:              "Name",
				Label:             "Name",
				Initial:           "a",
				TypeName:          "string",
				JSONName:          "name",
				ViewType:          "input",
				IsReference:       false,
				IsNested:          false,
				ReferenceName:     "",
				ReferenceJSONTags: []string{},
			},
			{
				Name:              "Age",
				Label:             "Age",
				Initial:           "a",
				TypeName:          "int",
				JSONName:          "age",
				ViewType:          "input",
				IsReference:       false,
				IsNested:          false,
				ReferenceName:     "",
				ReferenceJSONTags: []string{},
			},
		},
		HasReferences: false,
	}
}
	`))

	if err != nil {
		s.T().Fatal(err)
	}

	buf, err := s.gt.generate(content.TypePlain, t)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expectedBuffer, buf)
	}
}

func TestGenerate(t *testing.T) {
	suite.Run(t, new(GenerateTestSuite))
}
