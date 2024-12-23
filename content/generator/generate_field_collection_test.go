package generator

import (
	"encoding/json"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go/format"
	"testing"
)

type pageContentBlocks []content.FieldCollection

type ImageGallery struct {
	Headline   string `json:"headline"`
	Link       string `json:"link"`
	ButtonText string `json:"button_text"`
	Image      string `json:"image"`
}

type TextBlock struct {
	Text string `json:"text"`
}

type Link struct {
	ExternalUrl string `json:"external_url"`
	Label       string `json:"label"`
}

type imageAndTextBlock struct {
	Image         string `json:"image"`
	ImagePosition string `json:"image_position"`
	Content       string `json:"entities"`
	Link          Link   `json:"link"`
}

func (p *pageContentBlocks) Name() string {
	return "Page Content Blocks"
}

func (p *pageContentBlocks) Data() []content.FieldCollection {
	return *p
}

func (p *pageContentBlocks) Add(fieldCollection content.FieldCollection) {
	*p = append(*p, fieldCollection)
}

func (p *pageContentBlocks) Set(i int, fieldCollection content.FieldCollection) {
	data := p.Data()
	data[i] = fieldCollection
	*p = data
}

func (p *pageContentBlocks) SetData(data []content.FieldCollection) {
	*p = data
}

func (p *pageContentBlocks) AllowedTypes() map[string]content.Builder {
	return map[string]content.Builder{
		"ImageGallery": func() interface{} {
			return new(ImageGallery)
		},
		"ImageAndTextBlock": func() interface{} {
			return new(imageAndTextBlock)
		},
		"TextBlock": func() interface{} {
			return new(TextBlock)
		},
	}
}

func (p *pageContentBlocks) UnmarshalJSON(b []byte) error {
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

type GenerateFieldCollectionTestSuite struct {
	suite.Suite
	gt *contentGenerator
}

func (s *GenerateFieldCollectionTestSuite) SetupSuite() {
	var err error
	s.gt, err = setupGenerator(generator.Config{
		Target: generator.Target{
			Package: "entities",
		},
	}, content.Types{
		Definitions: map[string]generator.TypeDefinition{
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

func (s *GenerateFieldCollectionTestSuite) TestGenerateFieldCollectionField() {
	t := &generator.TypeDefinition{
		Name:  "Page",
		Label: "Page",
		Type:  generator.Content,
		Metadata: generator.Metadata{
			MethodReceiverName: "p",
		},
		Blocks: []generator.Block{
			{
				Type:     generator.Field,
				Name:     "Title",
				Label:    "Title",
				TypeName: "string",
				JSONName: "title",
				Definition: generator.BlockDefinition{
					Title: "Title",
					Type:  "string",
				},
			},
			{
				Type:     generator.Field,
				Name:     "URL",
				Label:    "URL",
				TypeName: "string",
				JSONName: "url",
				Definition: generator.BlockDefinition{
					Title: "URL",
					Type:  "string",
				},
			},
			{
				Type:          generator.Field,
				Name:          "ContentBlocks",
				Label:         "Content Blocks",
				TypeName:      "PageContentBlocks",
				ReferenceName: "PageContentBlocks",
				JSONName:      "content_blocks",
				Definition: generator.BlockDefinition{
					Title:       "ContentBlocks",
					Type:        "PageContentBlocks",
					IsReference: true,
				},
			},
		},
	}

	expectedBuffer, err := format.Source([]byte(`
package entities

import (
	"fmt"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/tokens"
)

type Page struct {
	item.Item

	Title         string             ` + "`json:\"title\"`" + `
	URL           string             ` + "`json:\"url\"`" + `
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
			View: editor.Input("Title", p, map[string]string{
				"label":       "Title",
				"type":        "text",
				"placeholder": "Enter the Title here",
			}, nil),
		},
		editor.Field{
			View: editor.Input("URL", p, map[string]string{
				"label":       "URL",
				"type":        "text",
				"placeholder": "Enter the URL here",
			}, nil),
		},
		editor.Field{
			View: editor.FieldCollection(
				"ContentBlocks",
				"Content Blocks",
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

func TestGenerateFieldCollection(t *testing.T) {
	suite.Run(t, new(GenerateFieldCollectionTestSuite))
}
