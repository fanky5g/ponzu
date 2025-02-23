package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeDefinition(t *testing.T) {
	tt := []struct {
		name                   string
		definitionType         Type
		tokens                 []string
		expectedTypeDefinition *TypeDefinition
	}{
		{
			name:           "PlainTypeDefinitionWithReferenceField",
			definitionType: Plain,
			// author name:string age:int image:@image
			tokens: []string{
				"author",
				"name:string",
				"age:int",
				"image:@image",
			},
			expectedTypeDefinition: &TypeDefinition{
				Name:  "Author",
				Label: "Author",
				Blocks: []Block{
					{
						Type:          Field,
						Name:          "Name",
						Label:         "Name",
						JSONName:      "name",
						TypeName:      "string",
						ReferenceName: "",
						Definition: BlockDefinition{
							Title:       "name",
							Type:        "string",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:          Field,
						Name:          "Age",
						Label:         "Age",
						JSONName:      "age",
						TypeName:      "int",
						ReferenceName: "",
						Definition: BlockDefinition{
							Title:       "age",
							Type:        "int",
							IsArray:     false,
							IsReference: false,
						},
					},
					{
						Type:          Field,
						Name:          "Image",
						Label:         "Image",
						JSONName:      "image",
						TypeName:      "string",
						ReferenceName: "Image",
						Definition: BlockDefinition{
							Title:       "image",
							Type:        "@image",
							IsArray:     false,
							IsReference: true,
						},
					},
				},
				Type: Plain,
				Metadata: Metadata{
					MethodReceiverName: "a",
				},
			},
		},
		{
			name:           "ParseContent",
			definitionType: Content,
			// blog title:string Author:string PostCategory:string content:string some_thing:int
			tokens: []string{
				"blog",
				"title:string",
				"Author:string",
				"PostCategory:string",
				"content:string",
				"some_thing:int",
			},
			expectedTypeDefinition: &TypeDefinition{
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
			},
		},
		{
			name:           "ParseFieldCollection",
			definitionType: FieldCollection,
			// page-content-blocks:"Page Content Blocks" ImageGallery:"Image Gallery" ImageAndTextBlock:"Image And Text Block"
			// TextBlock:"Text Block"
			tokens: []string{
				"page-content-blocks:Page Content Blocks",
				"ImageGallery:Image Gallery",
				"ImageAndTextBlock:Image And Text Block",
				"TextBlock:Text Block",
			},
			expectedTypeDefinition: &TypeDefinition{
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
			},
		},
		{
			name:           "ParseTypeWithReferenceField",
			definitionType: Content,
			// blog title:string author:@author category:string content:string
			tokens: []string{
				"blog",
				"title:string",
				"author:@author",
				"category:string",
				"content:string",
			},
			expectedTypeDefinition: &TypeDefinition{
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
						ReferenceName: "Author",
						Definition: BlockDefinition{
							Title:       "author",
							Type:        "@author",
							IsArray:     false,
							IsReference: true,
						},
					},
					{
						Type:          Field,
						Name:          "Category",
						Label:         "Category",
						JSONName:      "category",
						TypeName:      "string",
						ReferenceName: "",
						Definition: BlockDefinition{
							Title:       "category",
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
				},
				Type: Content,
				Metadata: Metadata{
					MethodReceiverName: "b",
				},
			},
		},
		{
			name:           "ParseTypeWithReferenceArrayField",
			definitionType: Content,
			// blog title:string authors:[]@author category:string content:string
			tokens: []string{
				"blog",
				"title:string",
				"authors:[]@author",
				"category:string",
				"content:string",
			},
			expectedTypeDefinition: &TypeDefinition{
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
						Name:          "Authors",
						Label:         "Authors",
						JSONName:      "authors",
						TypeName:      "[]string",
						ReferenceName: "Author",
						Definition: BlockDefinition{
							Title:       "authors",
							Type:        "[]@author",
							IsArray:     true,
							IsReference: true,
						},
					},
					{
						Type:          Field,
						Name:          "Category",
						Label:         "Category",
						JSONName:      "category",
						TypeName:      "string",
						ReferenceName: "",
						Definition: BlockDefinition{
							Title:       "category",
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
				},
				Type: Content,
				Metadata: Metadata{
					MethodReceiverName: "b",
				},
			},
		},
		{
			name: "ParseFieldCollectionAndReferenceField",
			// page author:@author content_blocks:@page_content_blocks
			tokens: []string{
				"page",
				"author:@author",
				"content_blocks:@page_content_blocks",
			},
			definitionType: Content,
			expectedTypeDefinition: &TypeDefinition{
				Name:  "Page",
				Label: "Page",
				Blocks: []Block{
					{
						Type:          Field,
						Name:          "Author",
						Label:         "Author",
						JSONName:      "author",
						TypeName:      "string",
						ReferenceName: "Author",
						Definition: BlockDefinition{
							Title:       "author",
							Type:        "@author",
							IsArray:     false,
							IsReference: true,
						},
					},
					{
						Type:          Field,
						Name:          "ContentBlocks",
						Label:         "ContentBlocks",
						TypeName:      "string",
						ReferenceName: "PageContentBlocks",
						JSONName:      "content_blocks",
						Definition: BlockDefinition{
							Title:       "content_blocks",
							Type:        "@page_content_blocks",
							IsReference: true,
						},
					},
				},
				Type: Content,
				Metadata: Metadata{
					MethodReceiverName: "p",
				},
			},
		},
		{
			name: "ParsePlainTypeWithTokens",
			// author name:string age:int gender:string:select@male~Male,female~Female,divers~Divers
			tokens: []string{
				"author",
				"name:string",
				"age:int",
				"gender:string:select@male~Male,female~Female,divers~Divers",
			},
			definitionType: Plain,
			expectedTypeDefinition: &TypeDefinition{
				Name:  "Author",
				Label: "Author",
				Blocks: []Block{
					{
						Type:     Field,
						Name:     "Name",
						Label:    "Name",
						JSONName: "name",
						TypeName: "string",
						Definition: BlockDefinition{
							Title: "name",
							Type:  "string",
						},
					},
					{
						Type:     Field,
						Name:     "Age",
						Label:    "Age",
						JSONName: "age",
						TypeName: "int",
						Definition: BlockDefinition{
							Title: "age",
							Type:  "int",
						},
					},
					{
						Type:     Field,
						Name:     "Gender",
						Label:    "Gender",
						JSONName: "gender",
						TypeName: "string",
						Definition: BlockDefinition{
							Title:  "gender",
							Type:   "string:select",
							Tokens: []string{"male:Male", "female:Female", "divers:Divers"},
						},
					},
				},
				Type: Plain,
				Metadata: Metadata{
					MethodReceiverName: "a",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gt, err := NewTypeDefinition(tc.definitionType, tc.tokens)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.expectedTypeDefinition, gt)
			}
		})
	}
}
