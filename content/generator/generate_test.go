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
	type image struct {
		item.Item
		Headline string `json:"headline"`
		URL      string `json:"url"`
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
			"Image": func() interface{} {
				return new(image)
			},
		},
		Definitions: map[string]generator.TypeDefinition{
			"Creator": {
				Name:  "Creator",
				Label: "Creator",
				Blocks: []generator.Block{
					{
						Type:          generator.Field,
						Name:          "Name",
						Label:         "Name",
						JSONName:      "name",
						TypeName:      "string",
						ReferenceName: "",
						Definition: generator.BlockDefinition{
							Title:       "name",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:          generator.Field,
						Name:          "Age",
						Label:         "Age",
						JSONName:      "age",
						TypeName:      "int",
						ReferenceName: "",
						Definition: generator.BlockDefinition{
							Title:       "age",
							Type:        "int",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:          generator.Field,
						Name:          "Image",
						Label:         "Image",
						JSONName:      "image",
						TypeName:      "string",
						ReferenceName: "Image",
						Definition: generator.BlockDefinition{
							Title:       "image",
							Type:        "@image",
							IsArray:     false,
							IsReference: true,
						},
					},
				},
				Type: generator.Plain,
				Metadata: generator.Metadata{
					MethodReceiverName: "a",
				},
			},
			"ImageGallery": {
				Type:     generator.Plain,
				Name:     "ImageGallery",
				Label:    "Image Gallery",
				Metadata: generator.Metadata{MethodReceiverName: "i"},
				Blocks: []generator.Block{
					{
						Type:              generator.Field,
						Name:              "Headline",
						Label:             "Headline",
						TypeName:          "string",
						JSONName:          "headline",
						ReferenceName:     "",
						ReferenceJSONTags: []string{},
						Definition: generator.BlockDefinition{
							Title:       "Headline",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:              generator.Field,
						Name:              "Link",
						Label:             "Link",
						TypeName:          "string",
						JSONName:          "link",
						ReferenceName:     "",
						ReferenceJSONTags: []string{},
						Definition: generator.BlockDefinition{
							Title:       "Link",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:              generator.Field,
						Name:              "ButtonText",
						Label:             "ButtonText",
						TypeName:          "string",
						JSONName:          "button_text",
						ReferenceName:     "",
						ReferenceJSONTags: []string{},
						Definition: generator.BlockDefinition{
							Title:       "ButtonText",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:              generator.Field,
						Name:              "Image",
						Label:             "Image",
						TypeName:          "string",
						JSONName:          "image",
						ReferenceName:     "",
						ReferenceJSONTags: []string{},
						Definition: generator.BlockDefinition{
							Title:       "Image",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
				},
			},
			"TextBlock": {
				Type:     generator.Plain,
				Name:     "TextBlock",
				Label:    "Text",
				Metadata: generator.Metadata{MethodReceiverName: "t"},
				Blocks: []generator.Block{
					{
						Type:              generator.Field,
						Name:              "Text",
						Label:             "Text",
						TypeName:          "string",
						JSONName:          "text",
						ReferenceName:     "",
						ReferenceJSONTags: []string{},
						Definition: generator.BlockDefinition{
							Title:       "Text",
							Type:        "string:richtext",
							IsArray:     false,
							IsReference: false,
						},
					},
				},
			},
			"ImageAndTextBlock": {
				Type:     generator.Plain,
				Name:     "ImageAndTextBlock",
				Label:    "Image and Text",
				Metadata: generator.Metadata{MethodReceiverName: "i"},
				Blocks: []generator.Block{
					{
						Type:              generator.Field,
						Name:              "Image",
						Label:             "Image",
						TypeName:          "string",
						JSONName:          "image",
						ReferenceName:     "",
						ReferenceJSONTags: []string{},
						Definition: generator.BlockDefinition{
							Title:       "Image",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:              generator.Field,
						Name:              "ImagePosition",
						Label:             "Image Position",
						TypeName:          "string",
						JSONName:          "image_position",
						ReferenceName:     "",
						ReferenceJSONTags: []string{},
						Definition: generator.BlockDefinition{
							Title:       "ImagePosition",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:              generator.Field,
						Name:              "Content",
						Label:             "Content",
						TypeName:          "string",
						JSONName:          "content",
						ReferenceName:     "",
						ReferenceJSONTags: []string{},
						Definition: generator.BlockDefinition{
							Title:       "Content",
							Type:        "string:richtext",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:              generator.Field,
						Name:              "Link",
						Label:             "Link",
						TypeName:          "Link",
						JSONName:          "link",
						ReferenceName:     "Link",
						ReferenceJSONTags: []string{},
						Definition: generator.BlockDefinition{
							Title:       "Link",
							Type:        "Link",
							IsArray:     false,
							IsReference: true,
						},
					},
				},
			},
			"Link": link,
		},
		FieldCollections: map[string]content.Builder{
			"PageContentBlocks": func() interface{} {
				return new(pageContentBlocks)
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

func (s *GenerateTestSuite) TestGenerateWithReferenceFieldAndFieldCollection() {
	typeDefinition := &generator.TypeDefinition{
		Name:  "Page",
		Label: "Page",
		Blocks: []generator.Block{
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
				Name:          "ContentBlocks",
				Label:         "ContentBlocks",
				TypeName:      "string",
				ReferenceName: "PageContentBlocks",
				JSONName:      "content_blocks",
				Definition: generator.BlockDefinition{
					Title:       "content_blocks",
					Type:        "@page_content_blocks",
					IsReference: true,
				},
			},
		},
		Type: generator.Content,
		Metadata: generator.Metadata{
			MethodReceiverName: "p",
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

type Page struct {
        item.Item

		Author string ` + "`json:\"author\" reference:\"Author\"`" + ` 
		ContentBlocks *PageContentBlocks ` + "`json:\"content_blocks\"`" + `
}

// MarshalEditor writes a buffer of views to edit a Page within the CMS
// and implements editor.Editable
func (p *Page) MarshalEditor(paths config.Paths) ([]byte, error) {
        view, err := editor.Form(p,
                paths,
				// Take note that the first argument to these Input-like functions
                // is the string version of each Page field, and must follow
                // this pattern for auto-decoding and auto-encoding reasons:
                editor.Field{
                        View: editor.ReferenceSelect(paths, "Author", p, map[string]string{
                                "label":       "Select Author",
                        },
						"Author",
                    ),
                },
				editor.Field{
					View: editor.FieldCollection(
					"ContentBlocks",
					"ContentBlocks",
					p,
					map[string]func(interface{}, *editor.FieldArgs, ...editor.Field) []byte{
						"ImageAndTextBlock": func(
							v interface{},
							args *editor.FieldArgs,
							injectFields ...editor.Field,
						) []byte {
							fields := append([]editor.Field{
								{
									View: editor.Input("Image", v, map[string]string{
										"label":       "Image",
										"type":        "text",
										"placeholder": "Enter the Image here",
									}, args),
								},
								{
									View: editor.Input("ImagePosition", v, map[string]string{
										"label":       "Image Position",
										"type":        "text",
										"placeholder": "Enter the Image Position here",
									}, args),
								},
								{
									View: editor.Richtext("Content", v, map[string]string{
										"label":       "Content",
										"placeholder": "Enter the Content here",
									}, args),
								},
								{
									View: editor.Nested("Link", v, args,
										editor.Field{
											View: editor.Input("Link.ExternalUrl", v, map[string]string{
												"label":       "ExternalUrl",
												"type":        "text",
												"placeholder": "Enter the ExternalUrl here",
											}, args),
										},
										editor.Field{
											View: editor.Input("Link.Label", v, map[string]string{
												"label":       "Label",
												"type":        "text",
												"placeholder": "Enter the Label here",
											}, args),
										},
									),
								},
							}, injectFields...)

							return editor.Nested("", v, args, fields...)
						},
						"ImageGallery": func(
							v interface{},
							args *editor.FieldArgs,
							injectFields ...editor.Field,
						) []byte {
							fields := append([]editor.Field{
								{
									View: editor.Input("Headline", v, map[string]string{
										"label":       "Headline",
										"type":        "text",
										"placeholder": "Enter the Headline here",
									}, args),
								},
								{
									View: editor.Input("Link", v, map[string]string{
										"label":       "Link",
										"type":        "text",
										"placeholder": "Enter the Link here",
									}, args),
								},
								{
									View: editor.Input("ButtonText", v, map[string]string{
										"label":       "ButtonText",
										"type":        "text",
										"placeholder": "Enter the ButtonText here",
									}, args),
								},
								{
									View: editor.Input("Image", v, map[string]string{
										"label":       "Image",
										"type":        "text",
										"placeholder": "Enter the Image here",
									}, args),
								},
							}, injectFields...)

							return editor.Nested("", v, args, fields...)
						},
						"TextBlock": func(
							v interface{},
							args *editor.FieldArgs,
							injectFields ...editor.Field,
						) []byte {
							fields := append([]editor.Field{
								{
									View: editor.Richtext("Text", v, map[string]string{
										"label":       "Text",
										"placeholder": "Enter the Text here",
									}, args),
								},
							}, injectFields...)

							return editor.Nested("", v, args, fields...)
						},
					}),
				},
        )

        if err != nil {
                return nil, fmt.Errorf("failed to render Page editor view: %s", err.Error())
        }

        return view, nil
}

func init() {
        Content["Page"] = func() interface{} { return new(Page) }
}

func (p *Page) EntityName() string {
        return "Page"
}

func (p *Page) GetTitle() string {
        return p.ID
}

func (p *Page) GetRepositoryToken() tokens.RepositoryToken {
        return "page"
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

func (s *GenerateTestSuite) TestGeneratePlainTypeWithReferenceField() {
	typeDefinition := &generator.TypeDefinition{
		Name:  "Author",
		Label: "Author",
		Blocks: []generator.Block{
			{
				Type:          generator.Field,
				Name:          "Name",
				Label:         "Name",
				JSONName:      "name",
				TypeName:      "string",
				ReferenceName: "",
				Definition: generator.BlockDefinition{
					Title:       "name",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:          generator.Field,
				Name:          "Age",
				Label:         "Age",
				JSONName:      "age",
				TypeName:      "int",
				ReferenceName: "",
				Definition: generator.BlockDefinition{
					Title:       "age",
					Type:        "int",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:          generator.Field,
				Name:          "Image",
				Label:         "Image",
				JSONName:      "image",
				TypeName:      "string",
				ReferenceName: "Image",
				Definition: generator.BlockDefinition{
					Title:       "image",
					Type:        "@image",
					IsArray:     false,
					IsReference: true,
				},
			},
		},
		Type: generator.Plain,
		Metadata: generator.Metadata{
			MethodReceiverName: "a",
		},
	}

	expectedBuffer, err := format.Source([]byte(`
package entities

import (
        "github.com/fanky5g/ponzu/generator"
)

type Author struct {
	Name string ` + "`json:\"name\"`" + `
	Age int ` + "`json:\"age\"`" + `
	Image string ` + "`json:\"image\" reference:\"Image\"`" + ` 
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
					Title:       "name",
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
					Title:       "age",
					Type:        "int",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:              generator.Field,
				Name:              "Image",
				Label:             "Image",
				TypeName:          "string",
				JSONName:          "image",
				ReferenceName:     "Image",
				ReferenceJSONTags: []string{},
				Definition: generator.BlockDefinition{
					Title:       "image",
					Type:        "@image",
					IsArray:     false,
					IsReference: true,
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

	err = s.gt.Generate(typeDefinition, w)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), string(expectedBuffer), string(w.buf))
	}
}

func (s *GenerateTestSuite) TestGenerateContentWithPlainTypeHavingReferenceField() {
	typeDefinition := &generator.TypeDefinition{
		Name:  "Story",
		Label: "Story",
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
				Name:          "Body",
				Label:         "Body",
				JSONName:      "body",
				TypeName:      "string",
				ReferenceName: "",
				Definition: generator.BlockDefinition{
					Title:       "body",
					Type:        "string:richtext",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:          generator.Field,
				Name:          "Author",
				Label:         "Author",
				JSONName:      "author",
				// TODO: Split references and nested Types
				TypeName:      "Creator",
				ReferenceName: "Creator",
				Definition: generator.BlockDefinition{
					Title:       "author",
					Type:        "Creator",
					IsArray:     false,
					IsReference: true,
				},
			},
		},
		Type: generator.Content,
		Metadata: generator.Metadata{
			MethodReceiverName: "s",
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

type Story struct {
	item.Item

	Title string ` + "`json:\"title\"`" + ` 
	Body string ` + "`json:\"body\"`" + ` 
	Author Creator ` + "`json:\"author\"`" + ` 
}

// MarshalEditor writes a buffer of views to edit a Story within the CMS
// and implements editor.Editable
func (s *Story) MarshalEditor(paths config.Paths) ([]byte, error) {
        view, err := editor.Form(s,
                paths,
				// Take note that the first argument to these Input-like functions
                // is the string version of each Story field, and must follow
                // this pattern for auto-decoding and auto-encoding reasons:
                editor.Field{
					View: editor.Input("Title", s, map[string]string{
						"label":       "Title",
						"type":        "text",
						"placeholder": "Enter the Title here",
					}, nil),
				},
				editor.Field{
					View: editor.Richtext("Body", s, map[string]string{
						"label":       "Body",
						"placeholder": "Enter the Body here",
					}, nil),
				},
				editor.Field{
					View: editor.Nested("Author", s, nil,
						editor.Field{
							View: editor.Input("Author.Name", s, map[string]string{
								"label":       "Name",
								"type":        "text",
								"placeholder": "Enter the Name here",
							}, nil),
						},
						editor.Field{
							View: editor.Input("Author.Age", s, map[string]string{
								"label":       "Age",
								"type":        "text",
								"placeholder": "Enter the Age here",
							}, nil),
						},
						editor.Field{
							View: editor.ReferenceSelect(paths, "Author.Image", s, map[string]string{
									"label": "Select Image",
								},
								"Image",
							),  
 						},
					),
				},
        )

        if err != nil {
                return nil, fmt.Errorf("failed to render Story editor view: %s", err.Error())
        }

        return view, nil
}

func init() {
        Content["Story"] = func() interface{} { return new(Story) }
}

func (s *Story) EntityName() string {
        return "Story"
}

func (s *Story) GetTitle() string {
        return s.ID
}

func (s *Story) GetRepositoryToken() tokens.RepositoryToken {
        return "story"
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
