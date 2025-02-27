package content

import (
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/item"
	"go/format"
	"strings"
	"testing"

	"github.com/fanky5g/ponzu/generator"
	"github.com/stretchr/testify/assert"
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

func TestGenerate(t *testing.T) {
	type author struct {
		item.Item
		Name string `json:"name"`
	}
	type image struct {
		item.Item
		Headline string `json:"headline"`
		URL      string `json:"url"`
	}

	contentTypes := content.Types{
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
						Type:          generator.Field,
						Name:          "Headline",
						Label:         "Headline",
						TypeName:      "string",
						JSONName:      "headline",
						ReferenceName: "",
						Definition: generator.BlockDefinition{
							Title:       "Headline",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:          generator.Field,
						Name:          "Link",
						Label:         "Link",
						TypeName:      "string",
						JSONName:      "link",
						ReferenceName: "",
						Definition: generator.BlockDefinition{
							Title:       "Link",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:          generator.Field,
						Name:          "ButtonText",
						Label:         "ButtonText",
						TypeName:      "string",
						JSONName:      "button_text",
						ReferenceName: "",
						Definition: generator.BlockDefinition{
							Title:       "ButtonText",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:          generator.Field,
						Name:          "Image",
						Label:         "Image",
						TypeName:      "string",
						JSONName:      "image",
						ReferenceName: "",
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
						Type:          generator.Field,
						Name:          "Text",
						Label:         "Text",
						TypeName:      "string",
						JSONName:      "text",
						ReferenceName: "",
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
						Type:          generator.Field,
						Name:          "Image",
						Label:         "Image",
						TypeName:      "string",
						JSONName:      "image",
						ReferenceName: "",
						Definition: generator.BlockDefinition{
							Title:       "Image",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:          generator.Field,
						Name:          "ImagePosition",
						Label:         "Image Position",
						TypeName:      "string",
						JSONName:      "image_position",
						ReferenceName: "",
						Definition: generator.BlockDefinition{
							Title:       "ImagePosition",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:          generator.Field,
						Name:          "Content",
						Label:         "Content",
						TypeName:      "string",
						JSONName:      "content",
						ReferenceName: "",
						Definition: generator.BlockDefinition{
							Title:       "Content",
							Type:        "string:richtext",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:          generator.Field,
						Name:          "Link",
						Label:         "Link",
						TypeName:      "Link",
						JSONName:      "link",
						ReferenceName: "Link",
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
			"ButtonLink": {
				Type:  generator.Plain,
				Name:  "ButtonLink",
				Label: "ButtonLink",
				Metadata: generator.Metadata{
					MethodReceiverName: "b",
				},
				Blocks: []generator.Block{
					{
						Type:          generator.Field,
						Name:          "Type",
						Label:         "Type",
						TypeName:      "string",
						JSONName:      "type",
						ReferenceName: "",
						Definition: generator.BlockDefinition{
							Title:       "type",
							Type:        "string:select",
							IsArray:     false,
							IsReference: false,
							Tokens: []string{
								"outlined:Outlined",
								"text:Text",
								"contained:Contained",
							},
						},
					},
					{
						Type:          generator.Field,
						Name:          "Text",
						Label:         "Text",
						TypeName:      "string",
						JSONName:      "text",
						ReferenceName: "",
						Definition: generator.BlockDefinition{
							Title:       "text",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
							Tokens:      []string{},
						},
					},
					{
						Type:          generator.Field,
						Name:          "Link",
						Label:         "Link",
						TypeName:      "string",
						JSONName:      "link",
						ReferenceName: "Link",
						Definition: generator.BlockDefinition{
							Title:       "link",
							Type:        "@link",
							IsArray:     false,
							IsReference: true,
							Tokens:      []string{},
						},
					},
				},
			},
		},
		FieldCollections: map[string]content.Builder{
			"PageContentBlocks": func() interface{} {
				return new(pageContentBlocks)
			},
		},
	}

	generatorConfig := generator.Config{
		Target: generator.Target{
			Path: generator.Path{
				Root: "",
				Base: "",
			},
			Package: "entities",
		},
	}

	gt, err := setupGenerator(generatorConfig, contentTypes)
	if err != nil {
		t.Fatal(err)
	}

	tt := []struct {
		name           string
		typeDefinition *generator.TypeDefinition
		expectedOutput string
	}{
		{
			name: "FieldCollection",
			typeDefinition: &generator.TypeDefinition{
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
			},
			expectedOutput: `
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
`,
		},
		{
			name: "ContentType",
			typeDefinition: &generator.TypeDefinition{
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
			},
			expectedOutput: `
package entities

import (
	"fmt"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
)

type Author struct {
	item.Item

	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
}

// MarshalEditor writes a buffer of views to edit a Author within the CMS
// and implements editor.Editable
func (a *Author) MarshalEditor(publicPath string) ([]byte, error) {
	view, err := editor.Form(a,
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

func (a *Author) GetRepositoryToken() string {
	return "author"
}`,
		},
		{
			name: "PlainType",
			// blog title:string author:@author category:string content:string
			typeDefinition: &generator.TypeDefinition{
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
			},
			expectedOutput: `
package entities

import (
        "fmt"
        "github.com/fanky5g/ponzu/content/editor"
        "github.com/fanky5g/ponzu/content/item"
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
func (b *Blog) MarshalEditor(publicPath string) ([]byte, error) {
        view, err := editor.Form(b,
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
                        View: editor.ReferenceSelect(publicPath, "Author", b, map[string]string{
                                "label":       "Select Author",
                        },
						nil,
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

func (b *Blog) GetRepositoryToken() string {
        return "blog"
}
	`,
		},
		{
			name: "ContentTypeWithReferenceField",
			typeDefinition: &generator.TypeDefinition{
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
			},
			expectedOutput: `
package entities

import (
        "fmt"
        "github.com/fanky5g/ponzu/content/editor"
        "github.com/fanky5g/ponzu/content/item"
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
func (b *Blog) MarshalEditor(publicPath string) ([]byte, error) {
        view, err := editor.Form(b,
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
                        View: editor.ReferenceSelect(publicPath, "Author", b, map[string]string{
                                "label":       "Select Author",
                        },
						nil,
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

func (b *Blog) GetRepositoryToken() string {
        return "blog"
}
	`,
		},
		{
			name: "UploadReferenceField",
			typeDefinition: &generator.TypeDefinition{
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
						Name:          "Image",
						Label:         "Image",
						JSONName:      "image",
						TypeName:      "string",
						ReferenceName: "Upload",
						Definition: generator.BlockDefinition{
							Title:       "image",
							Type:        "@upload",
							IsArray:     false,
							IsReference: true,
						},
					},
				},
				Type: generator.Content,
				Metadata: generator.Metadata{
					MethodReceiverName: "a",
				},
			},
			expectedOutput: `
package entities

import (
        "fmt"
        "github.com/fanky5g/ponzu/content/editor"
        "github.com/fanky5g/ponzu/content/item"
)

type Author struct {
        item.Item

		Name string ` + "`json:\"name\"`" + `
		Image string ` + "`json:\"image\" reference:\"Upload\"`" + ` 
}

// MarshalEditor writes a buffer of views to edit a Author within the CMS
// and implements editor.Editable
func (a *Author) MarshalEditor(publicPath string) ([]byte, error) {
        view, err := editor.Form(a,
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
                        View: editor.ReferenceSelect(publicPath, "Image", a, map[string]string{
                                "label":       "Select Image",
                        },
						nil,
						"Upload",
                    ),
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

func (a *Author) GetRepositoryToken() string {
        return "author"
}
	`,
		},
		{
			name: "ReferenceArrayField",
			typeDefinition: &generator.TypeDefinition{
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
			},
			expectedOutput: `
package entities

import (
        "fmt"
        "github.com/fanky5g/ponzu/content/editor"
        "github.com/fanky5g/ponzu/content/item"
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
func (b *Blog) MarshalEditor(publicPath string) ([]byte, error) {
        view, err := editor.Form(b,
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
                        View: editor.ReferenceSelectRepeater(publicPath, "Authors", b, map[string]string{
                                "label":       "Select Authors",
                        },
						nil,
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

func (b *Blog) GetRepositoryToken() string {
        return "blog"
}
	`,
		},
		{
			name: "ReferenceFieldAndFieldCollection",
			typeDefinition: &generator.TypeDefinition{
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
			},
			expectedOutput: `
package entities

import (
        "fmt"
        "github.com/fanky5g/ponzu/content/editor"
        "github.com/fanky5g/ponzu/content/item"
)

type Page struct {
        item.Item

		Author string ` + "`json:\"author\" reference:\"Author\"`" + ` 
		ContentBlocks *PageContentBlocks ` + "`json:\"content_blocks\"`" + `
}

// MarshalEditor writes a buffer of views to edit a Page within the CMS
// and implements editor.Editable
func (p *Page) MarshalEditor(publicPath string) ([]byte, error) {
        view, err := editor.Form(p,
				// Take note that the first argument to these Input-like functions
                // is the string version of each Page field, and must follow
                // this pattern for auto-decoding and auto-encoding reasons:
                editor.Field{
                        View: editor.ReferenceSelect(publicPath, "Author", p, map[string]string{
                                "label":       "Select Author",
                        },
						nil,
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

func (p *Page) GetRepositoryToken() string {
        return "page"
}
	`,
		},
		{
			name: "PlainTypeWithReferenceField",
			typeDefinition: &generator.TypeDefinition{
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
			},
			expectedOutput: `
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
				Definition: generator.BlockDefinition{
					Title:       "name",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
					Tokens: []string{},
				},
			},
			{
				Type:              generator.Field,
				Name:              "Age",
				Label:             "Age",
				TypeName:          "int",
				JSONName:          "age",
				ReferenceName:     "",
				Definition: generator.BlockDefinition{
					Title:       "age",
					Type:        "int",
					IsArray:     false,
					IsReference: false,
					Tokens: []string{},
				},
			},
			{
				Type:              generator.Field,
				Name:              "Image",
				Label:             "Image",
				TypeName:          "string",
				JSONName:          "image",
				ReferenceName:     "Image",
				Definition: generator.BlockDefinition{
					Title:       "image",
					Type:        "@image",
					IsArray:     false,
					IsReference: true,
					Tokens: []string{},
				},
			},
		},
	}
}
	`,
		},
		{
			name: "ContentWithPlainTypeHavingReferenceField",
			typeDefinition: &generator.TypeDefinition{
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
						Type:     generator.Field,
						Name:     "Author",
						Label:    "Author",
						JSONName: "author",
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
			},
			expectedOutput: `
package entities

import (
        "fmt"
        "github.com/fanky5g/ponzu/content/editor"
        "github.com/fanky5g/ponzu/content/item"
)

type Story struct {
	item.Item

	Title string ` + "`json:\"title\"`" + ` 
	Body string ` + "`json:\"body\"`" + ` 
	Author Creator ` + "`json:\"author\"`" + ` 
}

// MarshalEditor writes a buffer of views to edit a Story within the CMS
// and implements editor.Editable
func (s *Story) MarshalEditor(publicPath string) ([]byte, error) {
        view, err := editor.Form(s,
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
							View: editor.ReferenceSelect(publicPath, "Author.Image", s, map[string]string{
									"label": "Select Image",
								},
								nil,
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

func (s *Story) GetRepositoryToken() string {
        return "story"
}
`,
		},
		{
			name: "PlainTypeWithSelectTokens",
			typeDefinition: &generator.TypeDefinition{
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
						Type:     generator.Field,
						Name:     "Gender",
						Label:    "Gender",
						JSONName: "gender",
						TypeName: "string",
						Definition: generator.BlockDefinition{
							Title:  "gender",
							Type:   "string:select",
							Tokens: []string{"male:Male", "female:Female", "divers:Divers"},
						},
					},
				},
				Type: generator.Plain,
				Metadata: generator.Metadata{
					MethodReceiverName: "a",
				},
			},
			expectedOutput: `
package entities

import (
        "github.com/fanky5g/ponzu/generator"
)

type Author struct {
	Name string ` + "`json:\"name\"`" + `
	Age int ` + "`json:\"age\"`" + `
	Gender string ` + "`json:\"gender\"`" + ` 
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
				Definition: generator.BlockDefinition{
					Title:       "name",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
					Tokens: []string{},
				},
			},
			{
				Type:              generator.Field,
				Name:              "Age",
				Label:             "Age",
				TypeName:          "int",
				JSONName:          "age",
				ReferenceName:     "",
				Definition: generator.BlockDefinition{
					Title:       "age",
					Type:        "int",
					IsArray:     false,
					IsReference: false,
					Tokens: []string{},
				},
			},
			{
				Type:              generator.Field,
				Name:              "Gender",
				Label:             "Gender",
				TypeName:          "string",
				JSONName:          "gender",
				ReferenceName:     "",
				Definition: generator.BlockDefinition{
					Title:       "gender",
					Type:        "string:select",
					IsArray:     false,
					IsReference: false,
					Tokens: []string{
						"male:Male",
						"female:Female",
						"divers:Divers",
					},
				},
			},
		},
	}
}
	`,
		},
		{
			name: "ContentTypeWithSelectTokens",
			typeDefinition: &generator.TypeDefinition{
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
						Type:     generator.Field,
						Name:     "Gender",
						Label:    "Gender",
						JSONName: "gender",
						TypeName: "string",
						Definition: generator.BlockDefinition{
							Title:  "gender",
							Type:   "string:select",
							Tokens: []string{"male:Male", "female:Female", "divers:Divers"},
						},
					},
				},
				Type: generator.Content,
				Metadata: generator.Metadata{
					MethodReceiverName: "a",
				},
			},
			expectedOutput: `
package entities

import (
        "fmt"
        "github.com/fanky5g/ponzu/content/editor"
        "github.com/fanky5g/ponzu/content/item"
)

type Author struct {
        item.Item

		Name   string ` + "`json:\"name\"`" + `
		Age    int ` + "`json:\"age\"`" + `
		Gender string ` + "`json:\"gender\"`" + `
}

// MarshalEditor writes a buffer of views to edit a Author within the CMS
// and implements editor.Editable
func (a *Author) MarshalEditor(publicPath string) ([]byte, error) {
        view, err := editor.Form(a,
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
                editor.Field{
                        View: editor.Select("Gender", a, map[string]string{
                                "label":       "Select Gender",
                        }, []string{
							"male:Male",
							"female:Female",
							"divers:Divers",
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

func (a *Author) GetRepositoryToken() string {
        return "author"
}
	`,
		},
		{
			name: "ContentTypeWithNestedRepeatedField",
			typeDefinition: &generator.TypeDefinition{
				Name:  "Banner",
				Label: "Banner",
				Blocks: []generator.Block{
					{
						Type:          generator.Field,
						Name:          "Text",
						Label:         "Text",
						JSONName:      "text",
						TypeName:      "string",
						ReferenceName: "",
						Definition: generator.BlockDefinition{
							Title:       "text",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:          generator.Field,
						Name:          "Cta",
						Label:         "Cta",
						JSONName:      "cta",
						TypeName:      "[]string",
						ReferenceName: "ButtonLink",
						Definition: generator.BlockDefinition{
							Title:       "cta",
							Type:        "[]@button_link",
							IsArray:     true,
							IsReference: true,
						},
					},
				},
				Type: generator.Content,
				Metadata: generator.Metadata{
					MethodReceiverName: "a",
				},
			},
			expectedOutput: `
package entities

import (
	"fmt"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
)

type Banner struct {
	item.Item

	Text   string ` + "`json:\"text\"`" + `
	Cta  []ButtonLink ` + "`json:\"cta\"`" + `
}

// MarshalEditor writes a buffer of views to edit a Banner within the CMS
// and implements editor.Editable
func (a *Banner) MarshalEditor(publicPath string) ([]byte, error) {
	view, err := editor.Form(a,
		// Take note that the first argument to these Input-like functions
		// is the string version of each Banner field, and must follow
		// this pattern for auto-decoding and auto-encoding reasons:
		editor.Field{
			View: editor.Input("Text", a, map[string]string{
				"label":       "Text",
				"type":        "text",
				"placeholder": "Enter the Text here",
			}, nil),
		},
		editor.Field{
			View: editor.NestedRepeater("Cta", a, func(v interface{}, args *editor.FieldArgs) (string, []editor.Field) {
				return "ButtonLink", []editor.Field{
					{
						View: editor.Select("Cta.Type", a, map[string]string{
							"label": "Select Type",
						}, []string{
							"outlined:Outlined",
							"text:Text",
							"contained:Contained",
						}, args),
					},
					{
						View: editor.Input("Cta.Text", a, map[string]string{
							"label":       "Text",
							"type":        "text",
							"placeholder": "Enter the Text here",
						}, args),
					},
					{
						View: editor.Nested("Cta.Link", a, args,
							editor.Field{
								View: editor.Input("Cta.Link.ExternalUrl", a, map[string]string{
									"label":       "ExternalUrl",
									"type":        "text",
									"placeholder": "Enter the ExternalUrl here",
								}, args),
							},
							editor.Field{
								View: editor.Input("Cta.Link.Label", a, map[string]string{
									"label":       "Label",
									"type":        "text",
									"placeholder": "Enter the Label here",
								}, args),
							},
						),
					},
				}
			},
			),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to render Banner editor view: %s", err.Error())
	}

	return view, nil
}

func init() {
	Content["Banner"] = func() interface{} { return new(Banner) }
}

func (a *Banner) EntityName() string {
	return "Banner"
}

func (a *Banner) GetTitle() string {
	return a.ID
}

func (a *Banner) GetRepositoryToken() string {
	return "banner"
}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var expectedBuffer []byte
			expectedBuffer, err = format.Source([]byte(tc.expectedOutput))

			if err != nil {
				t.Fatal(err)
			}

			w := new(testWriter)

			err = gt.Generate(tc.typeDefinition, w)
			if assert.NoError(t, err) {
				assert.Equal(t, expectedBuffer, w.buf)
			}
		})
	}
}
