package generator

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go/format"
	"testing"
)

var (
	link = generator.TypeDefinition{
		Type:  generator.Plain,
		Name:  "Link",
		Label: "Link",
		Metadata: generator.Metadata{
			MethodReceiverName: "l",
		},
		Blocks: []generator.Block{
			{
				Type:              generator.Field,
				Name:              "ExternalUrl",
				Label:             "ExternalUrl",
				TypeName:          "string",
				JSONName:          "external_url",
				ReferenceName:     "",
				ReferenceJSONTags: []string{},
				Definition: generator.BlockDefinition{
					Title:       "ExternalUrl",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:              generator.Field,
				Name:              "Label",
				Label:             "Label",
				TypeName:          "l",
				JSONName:          "label",
				ReferenceName:     "",
				ReferenceJSONTags: []string{},
				Definition: generator.BlockDefinition{
					Title:       "Label",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
		},
	}

	author = generator.TypeDefinition{
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
			{
				Type:          generator.Field,
				Name:          "Website",
				Label:         "Website",
				TypeName:      "Link",
				JSONName:      "website",
				ReferenceName: "Link",
				Definition: generator.BlockDefinition{
					Title:       "Website",
					Type:        "Link",
					IsArray:     false,
					IsReference: true,
				},
			},
		},
	}
)

type GenerateNestedTestSuite struct {
	suite.Suite
	gt *contentGenerator
}

func (s *GenerateNestedTestSuite) SetupSuite() {
	var err error
	s.gt, err = setupGenerator(generator.Config{
		Target: generator.Target{
			Path: generator.Path{
				Root: "",
				Base: "",
			},
			Package: "entities",
		},
	}, content.Types{
		Definitions: map[string]generator.TypeDefinition{
			"Link":   link,
			"Author": author,
		},
	})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *GenerateNestedTestSuite) TestNestedFieldInPlainType() {
	expectedBuffer, err := format.Source([]byte(`
package entities

import (
	"github.com/fanky5g/ponzu/generator"
)

type Author struct {
	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
	Website  Link    ` + "`json:\"website\"`" + `
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
			{
				Type:              generator.Field,
				Name:              "Website",
				Label:             "Website",
				TypeName:          "Link",
				JSONName:          "website",
				ReferenceName:     "Link",
				ReferenceJSONTags: []string{},
				Definition: generator.BlockDefinition{
					Title:       "Website",
					Type:        "Link",
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

	err = s.gt.Generate(&author, w)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expectedBuffer, w.buf)
	}
}

func (s *GenerateNestedTestSuite) TestNestedFieldInContentType() {
	t := author
	t.Type = generator.Content
	expectedBuffer, err := format.Source([]byte(`
package entities

import (
	"fmt"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/tokens"
)

type Author struct {
	item.Item

	Name     string  ` + "`json:\"name\"`" + `
	Age      int     ` + "`json:\"age\"`" + `
	Website  Link    ` + "`json:\"website\"`" + `
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
			View: editor.Nested("Website", a, nil,
				editor.Field{
					View: editor.Input("Website.ExternalUrl", a, map[string]string{
						"label":       "ExternalUrl",
						"type":        "text",
						"placeholder": "Enter the ExternalUrl here",
					}, nil),
				},
				editor.Field{
					View: editor.Input("Website.Label", a, map[string]string{
						"label":       "Label",
						"type":        "text",
						"placeholder": "Enter the Label here",
					}, nil),
				},
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

func (a *Author) GetRepositoryToken() tokens.RepositoryToken {
	return "author"
}
`))

	if err != nil {
		s.T().Fatal(err)
	}

	w := new(testWriter)

	err = s.gt.Generate(&t, w)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expectedBuffer, w.buf)
	}
}

func (s *GenerateNestedTestSuite) TestDoublyNestedFieldInContentType() {
	t := generator.TypeDefinition{
		Type:  generator.Content,
		Name:  "Review",
		Label: "Review",
		Blocks: []generator.Block{
			{
				Type:     generator.Field,
				Name:     "Title",
				Label:    "Title",
				TypeName: "string",
				JSONName: "title",
				Definition: generator.BlockDefinition{
					Title:       "Title",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:     generator.Field,
				Name:     "Body",
				Label:    "Body",
				TypeName: "string",
				JSONName: "body",
				Definition: generator.BlockDefinition{
					Title:       "Body",
					Type:        "string:richtext",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:     generator.Field,
				Name:     "Rating",
				Label:    "Rating",
				TypeName: "int",
				JSONName: "rating",
				Definition: generator.BlockDefinition{
					Title:       "Rating",
					Type:        "int",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:     generator.Field,
				Name:     "Tags",
				Label:    "Tags",
				TypeName: "[]string",
				JSONName: "tags",
				Definition: generator.BlockDefinition{
					Title:       "Tags",
					Type:        "[]string",
					IsArray:     true,
					IsReference: false,
				},
			},
			{
				Type:          generator.Field,
				Name:          "Author",
				Label:         "Author",
				TypeName:      "Author",
				JSONName:      "author",
				ReferenceName: "Author",
				Definition: generator.BlockDefinition{
					Title:       "Author",
					Type:        "Author",
					IsArray:     false,
					IsReference: true,
				},
			},
		},
		Metadata: generator.Metadata{
			MethodReceiverName: "r",
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

type Review struct {
	item.Item

	Title       string    ` + "`json:\"title\"`" + `
	Body        string    ` + "`json:\"body\"`" + `
	Rating      int       ` + "`json:\"rating\"`" + `
	Tags        []string  ` + "`json:\"tags\"`" + `
	Author      Author    ` + "`json:\"author\"`" + `
}

// MarshalEditor writes a buffer of views to edit a Review within the CMS
// and implements editor.Editable
func (r *Review) MarshalEditor(publicPath string) ([]byte, error) {
	view, err := editor.Form(r,
		// Take note that the first argument to these Input-like functions
		// is the string version of each Review field, and must follow
		// this pattern for auto-decoding and auto-encoding reasons:
		editor.Field{
			View: editor.Input("Title", r, map[string]string{
				"label":       "Title",
				"type":        "text",
				"placeholder": "Enter the Title here",
			}, nil),
		},
		editor.Field{
			View: editor.Richtext("Body", r, map[string]string{
				"label":       "Body",
				"placeholder": "Enter the Body here",
			}, nil),
		},
		editor.Field{
			View: editor.Input("Rating", r, map[string]string{
				"label":       "Rating",
				"type":        "text",
				"placeholder": "Enter the Rating here",
			}, nil),
		},
		editor.Field{
			View: editor.InputRepeater("Tags", r, map[string]string{
				"label":       "Tags",
				"type":        "text",
				"placeholder": "Enter the Tags here",
			}),
		},
		editor.Field{
			View: editor.Nested("Author", r, nil,
				editor.Field{
					View: editor.Input("Author.Name", r, map[string]string{
						"label":       "Name",
						"type":        "text",
						"placeholder": "Enter the Name here",
					}, nil),
				},
				editor.Field{
					View: editor.Input("Author.Age", r, map[string]string{
						"label":       "Age",
						"type":        "text",
						"placeholder": "Enter the Age here",
					}, nil),
				},
				editor.Field{
					View: editor.Nested("Author.Website", r, nil,
						editor.Field{
							View: editor.Input("Author.Website.ExternalUrl", r, map[string]string{
								"label":       "ExternalUrl",
								"type":        "text",
								"placeholder": "Enter the ExternalUrl here",
							}, nil),
						},
						editor.Field{
							View: editor.Input("Author.Website.Label", r, map[string]string{
								"label":       "Label",
								"type":        "text",
								"placeholder": "Enter the Label here",
							}, nil),
						},
					),
				},
			),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to render Review editor view: %s", err.Error())
	}

	return view, nil
}

func init() {
	Content["Review"] = func() interface{} { return new(Review) }
}

func (r *Review) EntityName() string {
	return "Review"
}

func (r *Review) GetTitle() string {
	return r.ID
}

func (r *Review) GetRepositoryToken() tokens.RepositoryToken {
	return "review"
}
`))

	if err != nil {
		s.T().Fatal(err)
	}

	w := new(testWriter)

	err = s.gt.Generate(&t, w)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expectedBuffer, w.buf)
	}
}

func TestGenerateNested(t *testing.T) {
	suite.Run(t, new(GenerateNestedTestSuite))
}
