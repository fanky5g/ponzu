package generator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TypeDefinitionTestSuite struct {
	suite.Suite
}

func (s *TypeDefinitionTestSuite) TestParseContent() {
	// blog title:string Author:string PostCategory:string content:string some_thing:int
	args := []string{
		"blog",
		"title:string",
		"Author:string",
		"PostCategory:string",
		"content:string",
		"some_thing:int",
	}

	definition, err := NewTypeDefinition(Content, args)
	if err != nil {
		s.T().Errorf("Failed: %s", err.Error())
	}

	expectedTypeDefinition := &TypeDefinition{
		Name:  "Blog",
		Label: "Blog",
		Blocks: []Block{
			{
				Type:          Field,
				Name:          "Title",
				Label:         "Title",
				JSONName:      "title",
				TypeName:      "string",
				ReferenceName: "",
				Definition: BlockDefinition{
					Title:       "title",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:          Field,
				Name:          "Author",
				Label:         "Author",
				JSONName:      "author",
				TypeName:      "string",
				ReferenceName: "",
				Definition: BlockDefinition{
					Title:       "Author",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:          Field,
				Name:          "PostCategory",
				Label:         "PostCategory",
				JSONName:      "post_category",
				TypeName:      "string",
				ReferenceName: "",
				Definition: BlockDefinition{
					Title:       "PostCategory",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:          Field,
				Name:          "Content",
				Label:         "Content",
				JSONName:      "content",
				TypeName:      "string",
				ReferenceName: "",
				Definition: BlockDefinition{
					Title:       "content",
					Type:        "string",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:          Field,
				Name:          "SomeThing",
				Label:         "SomeThing",
				JSONName:      "some_thing",
				TypeName:      "int",
				ReferenceName: "",
				Definition: BlockDefinition{
					Title:       "some_thing",
					Type:        "int",
					IsArray:     false,
					IsReference: false,
				},
			},
		},
		Type: Content,
		Metadata: Metadata{
			MethodReceiverName: "b",
		},
	}

	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), definition, expectedTypeDefinition)
	}
}

func (s *TypeDefinitionTestSuite) TestParseFieldCollection() {
	// page-content-blocks:"Page Content Blocks" ImageGallery:"Image Gallery" ImageAndTextBlock:"Image And Text Block"
	// TextBlock:"Text Block"
	args := []string{
		"page-content-blocks:Page Content Blocks",
		"ImageGallery:Image Gallery",
		"ImageAndTextBlock:Image And Text Block",
		"TextBlock:Text Block",
	}

	gt, err := NewTypeDefinition(FieldCollection, args)
	expectedType := &TypeDefinition{
		Type:     FieldCollection,
		Name:     "PageContentBlocks",
		Label:    "Page Content Blocks",
		Metadata: Metadata{MethodReceiverName: "p"},
		Blocks: []Block{
			{
				Type:     ContentBlock,
				Name:     "ImageGallery",
				TypeName: "ImageGallery",
				Label:    "Image Gallery",
				JSONName: "image_gallery",
				Definition: BlockDefinition{
					Title:       "ImageGallery",
					Type:        "ImageGallery",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:     ContentBlock,
				Name:     "ImageAndTextBlock",
				TypeName: "ImageAndTextBlock",
				Label:    "Image And Text Block",
				JSONName: "image_and_text_block",
				Definition: BlockDefinition{
					Title:       "ImageAndTextBlock",
					Type:        "ImageAndTextBlock",
					IsArray:     false,
					IsReference: false,
				},
			},
			{
				Type:     ContentBlock,
				Name:     "TextBlock",
				TypeName: "TextBlock",
				Label:    "Text Block",
				JSONName: "text_block",
				Definition: BlockDefinition{
					Title:       "TextBlock",
					Type:        "TextBlock",
					IsArray:     false,
					IsReference: false,
				},
			},
		},
	}

	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expectedType, gt)
	}
}

func TestTypeDefinition(t *testing.T) {
	suite.Run(t, new(TypeDefinitionTestSuite))
}
