package generator

import (
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go/format"
	"strings"
	"testing"
)

type GenerateTestSuite struct {
	suite.Suite
	gt *contentGenerator
}

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

func (s *GenerateTestSuite) SetupSuite() {
	var err error
	s.gt, err = setupGenerator(generator.Config{
		Target: generator.Target{
			Path: generator.Path{
				Root: "",
				Base: "",
			},
			Package: "entities",
		},
	}, content.Types{})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *GenerateTestSuite) TestGenerateFieldCollection() {
	t := &generator.TypeDefinition{
		Name:     "PageContentBlocks",
		Label:    "Page Content Blocks",
		Type:     generator.FieldCollection,
		Metadata: generator.Metadata{MethodReceiverName: "p"},
		Blocks: []generator.Block{
			{
				Type:     generator.ContentBlock,
				TypeName: "ImageGallery",
				Label:    "Image Gallery",
			},
			{
				Type:     generator.ContentBlock,
				TypeName: "ImageAndTextBlock",
				Label:    "Image And Text Block",
			},
			{
				Type:     generator.ContentBlock,
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

	w := new(testWriter)

	err = s.gt.Generate(t, w)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expectedBuffer, w.buf)
	}
}

func (s *GenerateTestSuite) TestGenerateContentType() {
	t := &generator.TypeDefinition{
		Type:  generator.Content,
		Name:  "Author",
		Label: "Author",
		Metadata: generator.Metadata{
			MethodReceiverName: "a",
		},
		Blocks: []generator.Block{
			{
				Type:     generator.Field,
				Name:     "Name",
				Label:    "Name",
				TypeName: "string",
				JSONName: "name",
			},
			{
				Type:     generator.Field,
				Name:     "Age",
				Label:    "Age",
				TypeName: "int",
				JSONName: "age",
			},
		},
	}

	expectedBuffer, err := format.Source([]byte(`
package entities

import (
	"fmt"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/tokens"
	"reflect"
)

type Author struct {
	item.Item

	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
}

// MarshalEditor writes a buffer of views to edit a Author within the CMS
// and implements editor.Editable
func (a *Author) MarshalEditor(paths config.Paths) ([]byte, error) {
	view, err := editor.Form(a,
		paths,
		// Take note that the first argument to these Input-like functions
		// is the string version of each Author field, and must follow
		// this pattern for auto-decoding and auto-encoding reasons:
		editor.Field{
			View: editor.Input("Name", a, map[string]string{
				"label":       "Name",
				"type":        "text",
				"placeholder": "Enter the Name here",
			}, nil),
		},
		editor.Field{
			View: editor.Input("Age", a, map[string]string{
				"label":       "Age",
				"type":        "text",
				"placeholder": "Enter the Age here",
			}, nil),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to render Author editor view: %s", err.Error())
	}

	return view, nil
}

func init() {
	Content["Author"] = func() interface{} { return new(Author) }
}

// IndexContent determines if Author should be indexed for searching
func (a *Author) IndexContent() bool {
	return false
}

// GetSearchableAttributes defines fields of Author that should be indexed
func (a *Author) GetSearchableAttributes() map[string]reflect.Type {
	searchableAttributes := make(map[string]reflect.Type)
	idField := "ID"
	v := reflect.Indirect(reflect.ValueOf(a))
	searchableAttributes[idField] = v.FieldByName(idField).Type()

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		if fieldName != idField && field.Kind() == reflect.String {
			searchableAttributes[fieldName] = field.Type()
		}
	}

	return searchableAttributes
}

func (a *Author) GetTitle() string {
	return a.ID
}

func (a *Author) GetRepositoryToken() tokens.RepositoryToken {
	return "author"
}`))

	if err != nil {
		s.T().Fatal(err)
	}

	w := new(testWriter)

	err = s.gt.Generate(t, w)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expectedBuffer, w.buf)
	}
}

func (s *GenerateTestSuite) TestGeneratePlainType() {
	t := &generator.TypeDefinition{
		Type:  generator.Plain,
		Name:  "Author",
		Label: "Author",
		Metadata: generator.Metadata{
			MethodReceiverName: "a",
		},
		Blocks: []generator.Block{
			{
				Type:     generator.Field,
				Name:     "Name",
				Label:    "Name",
				TypeName: "string",
				JSONName: "name",
				Definition: generator.BlockDefinition{
					Title:       "Name",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:     generator.Field,
				Name:     "Age",
				Label:    "Age",
				TypeName: "int",
				JSONName: "age",
				Definition: generator.BlockDefinition{
					Title:       "Age",
					Type:        "int",
					IsArray:     false,
					IsReference: false,
				},
			},
		},
	}

	expectedBuffer, err := format.Source([]byte(`
package entities

import (
	"github.com/fanky5g/ponzu/generator"
)

type Author struct {
	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
}

func init() {
	Definitions["Author"] = generator.TypeDefinition{
		Type:  generator.Plain,
		Name:  "Author",
		Label: "Author",
		Metadata: generator.Metadata{
			MethodReceiverName: "a",
		},
		Blocks: []generator.Block{
			{
				Type:              generator.Field,
				Name:              "Name",
				Label:             "Name",
				TypeName:          "string",
				JSONName:          "name",
				ReferenceName:     "",
				ReferenceJSONTags: []string{},
				Definition: generator.BlockDefinition{
					Title:       "Name",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:              generator.Field,
				Name:              "Age",
				Label:             "Age",
				TypeName:          "int",
				JSONName:          "age",
				ReferenceName:     "",
				ReferenceJSONTags: []string{},
				Definition: generator.BlockDefinition{
					Title:       "Age",
					Type:        "int",
					IsArray:     false,
					IsReference: false,
				},
			},
		},
	}
}
	`))

	if err != nil {
		s.T().Fatal(err)
	}

	w := new(testWriter)

	err = s.gt.Generate(t, w)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expectedBuffer, w.buf)
	}
}

func TestGenerate(t *testing.T) {
	suite.Run(t, new(GenerateTestSuite))
}
