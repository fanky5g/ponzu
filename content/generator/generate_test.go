package generator

import (
	"fmt"
	"go/format"
	"strings"
	"testing"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	type author struct {
		item.Item
		Name string `json:"name"`
	}
	s.gt, err = setupGenerator(generator.Config{
		Target: generator.Target{
			Path: generator.Path{
				Root: "",
				Base: "",
			},
			Package: "entities",
		},
	}, content.Types{
		Content: map[string]content.Builder{
			"Author": func() interface{} {
				return new(author)
			},
		},
	})
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

func (a *Author) EntityName() string {
	return "Author"
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

func (s *GenerateTestSuite) TestGenerateWithReferenceField() {
	typeDefinition := &generator.TypeDefinition{
		Name:  "Blog",
		Label: "Blog",
		Blocks: []generator.Block{
			{
				Type:          generator.Field,
				Name:          "Title",
				Label:         "Title",
				JSONName:      "title",
				TypeName:      "string",
				ReferenceName: "",
				Definition: generator.BlockDefinition{
					Title:       "title",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:          generator.Field,
				Name:          "Author",
				Label:         "Author",
				JSONName:      "author",
				TypeName:      "string",
				ReferenceName: "Author",
				Definition: generator.BlockDefinition{
					Title:       "author",
					Type:        "@author",
					IsArray:     false,
					IsReference: true,
				},
			},
			{
				Type:          generator.Field,
				Name:          "Category",
				Label:         "Category",
				JSONName:      "category",
				TypeName:      "string",
				ReferenceName: "",
				Definition: generator.BlockDefinition{
					Title:       "category",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:          generator.Field,
				Name:          "Content",
				Label:         "Content",
				JSONName:      "content",
				TypeName:      "string",
				ReferenceName: "",
				Definition: generator.BlockDefinition{
					Title:       "content",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
		},
		Type: generator.Content,
		Metadata: generator.Metadata{
			MethodReceiverName: "b",
		},
	}

	// blog title:string author:@author category:string content:string
	expectedBuffer, err := format.Source([]byte(`
package entities

import (
        "fmt"
        "github.com/fanky5g/ponzu/config"
        "github.com/fanky5g/ponzu/content/editor"
        "github.com/fanky5g/ponzu/content/item"
        "github.com/fanky5g/ponzu/tokens"
)

type Blog struct {
        item.Item

		Title string ` + "`json:\"title\"`" + `
		Author string ` + "`json:\"author\" reference:\"Author\"`" + ` 
		Category string ` + "`json:\"category\"`" + `
		Content string ` + "`json:\"content\"`" + `
}

// MarshalEditor writes a buffer of views to edit a Blog within the CMS
// and implements editor.Editable
func (b *Blog) MarshalEditor(paths config.Paths) ([]byte, error) {
        view, err := editor.Form(b,
                paths,
                // Take note that the first argument to these Input-like functions
                // is the string version of each Blog field, and must follow
                // this pattern for auto-decoding and auto-encoding reasons:
                editor.Field{
                        View: editor.Input("Title", b, map[string]string{
                                "label":       "Title",
                                "type":        "text",
                                "placeholder": "Enter the Title here",
                        }, nil),
                },
                editor.Field{
                        View: editor.ReferenceSelect(paths, "Author", b, map[string]string{
                                "label":       "Select Author",
                        },
						"Author",
						` + "`Author: {{ .id }}`" + `,
                    ),
                },
                editor.Field{
                        View: editor.Input("Category", b, map[string]string{
                                "label":       "Category",
                                "type":        "text",
                                "placeholder": "Enter the Category here",
                        }, nil),
                },
                editor.Field{
                        View: editor.Input("Content", b, map[string]string{
                                "label":       "Content",
                                "type":        "text",
                                "placeholder": "Enter the Content here",
                        }, nil),
                },
        )

        if err != nil {
                return nil, fmt.Errorf("failed to render Blog editor view: %s", err.Error())
        }

        return view, nil
}

func init() {
        Content["Blog"] = func() interface{} { return new(Blog) }
}

func (b *Blog) EntityName() string {
        return "Blog"
}

func (b *Blog) GetTitle() string {
        return b.ID
}

func (b *Blog) GetRepositoryToken() tokens.RepositoryToken {
        return "blog"
}
	`))

	if err != nil {
		s.T().Fatal(err)
	}

	w := new(testWriter)

	err = s.gt.Generate(typeDefinition, w)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), string(expectedBuffer), string(w.buf))
	}

}

func (s *GenerateTestSuite) TestGenerateWithReferenceArrayField() {
	typeDefinition := &generator.TypeDefinition{
		Name:  "Blog",
		Label: "Blog",
		Blocks: []generator.Block{
			{
				Type:          generator.Field,
				Name:          "Title",
				Label:         "Title",
				JSONName:      "title",
				TypeName:      "string",
				ReferenceName: "",
				Definition: generator.BlockDefinition{
					Title:       "title",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:          generator.Field,
				Name:          "Authors",
				Label:         "Authors",
				JSONName:      "authors",
				TypeName:      "[]string",
				ReferenceName: "Author",
				Definition: generator.BlockDefinition{
					Title:       "authors",
					Type:        "[]@author",
					IsArray:     true,
					IsReference: true,
				},
			},
			{
				Type:          generator.Field,
				Name:          "Category",
				Label:         "Category",
				JSONName:      "category",
				TypeName:      "string",
				ReferenceName: "",
				Definition: generator.BlockDefinition{
					Title:       "category",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:          generator.Field,
				Name:          "Content",
				Label:         "Content",
				JSONName:      "content",
				TypeName:      "string",
				ReferenceName: "",
				Definition: generator.BlockDefinition{
					Title:       "content",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
		},
		Type: generator.Content,
		Metadata: generator.Metadata{
			MethodReceiverName: "b",
		},
	}

	// blog title:string author:@author category:string content:string
	expectedBuffer, err := format.Source([]byte(`
package entities

import (
        "fmt"
        "github.com/fanky5g/ponzu/config"
        "github.com/fanky5g/ponzu/content/editor"
        "github.com/fanky5g/ponzu/content/item"
        "github.com/fanky5g/ponzu/tokens"
)

type Blog struct {
        item.Item

		Title string ` + "`json:\"title\"`" + `
		Authors []string ` + "`json:\"authors\" reference:\"Author\"`" + ` 
		Category string ` + "`json:\"category\"`" + `
		Content string ` + "`json:\"content\"`" + `
}

// MarshalEditor writes a buffer of views to edit a Blog within the CMS
// and implements editor.Editable
func (b *Blog) MarshalEditor(paths config.Paths) ([]byte, error) {
        view, err := editor.Form(b,
                paths,
                // Take note that the first argument to these Input-like functions
                // is the string version of each Blog field, and must follow
                // this pattern for auto-decoding and auto-encoding reasons:
                editor.Field{
                        View: editor.Input("Title", b, map[string]string{
                                "label":       "Title",
                                "type":        "text",
                                "placeholder": "Enter the Title here",
                        }, nil),
                },
                editor.Field{
                        View: editor.ReferenceSelectRepeater(paths, "Authors", b, map[string]string{
                                "label":       "Authors",
                        },
						"Author",
						` + "`Author: {{ .id }}`" + `,
                    ),
                },
                editor.Field{
                        View: editor.Input("Category", b, map[string]string{
                                "label":       "Category",
                                "type":        "text",
                                "placeholder": "Enter the Category here",
                        }, nil),
                },
                editor.Field{
                        View: editor.Input("Content", b, map[string]string{
                                "label":       "Content",
                                "type":        "text",
                                "placeholder": "Enter the Content here",
                        }, nil),
                },
        )

        if err != nil {
                return nil, fmt.Errorf("failed to render Blog editor view: %s", err.Error())
        }

        return view, nil
}

func init() {
        Content["Blog"] = func() interface{} { return new(Blog) }
}

func (b *Blog) EntityName() string {
        return "Blog"
}

func (b *Blog) GetTitle() string {
        return b.ID
}

func (b *Blog) GetRepositoryToken() tokens.RepositoryToken {
        return "blog"
}
	`))

	if err != nil {
		s.T().Fatal(err)
	}

	w := new(testWriter)

	err = s.gt.Generate(typeDefinition, w)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), string(expectedBuffer), string(w.buf))
	}

}

func TestGenerate(t *testing.T) {
	suite.Run(t, new(GenerateTestSuite))
}
