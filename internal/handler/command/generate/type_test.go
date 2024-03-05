package generate

import (
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/domain/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ParseTypeTestSuite struct {
	suite.Suite
}

func (suite ParseTypeTestSuite) TestTypeName() {
	// blog title:string Author:string PostCategory:string content:string some_thing:int
	args := []string{
		"blog", "title:string", "Author:string",
		"PostCategory:string", "content:string",
		"some_thing:int", "Some_otherThing:float64",
	}

	gt, err := parseType(enum.TypeContent, args)
	if err != nil {
		suite.T().Errorf("Failed: %s", err.Error())
	}

	if gt.Name != "Blog" {
		suite.T().Errorf("Expected %s, got: %s", "Blog", gt.Name)
	}
}

func (suite ParseTypeTestSuite) TestParsing() {
	// page-content-blocks:"Page Content Blocks" ImageGallery:"Image Gallery" ImageAndTextBlock:"Image And Text Block"
	// TextBlock:"Text Block"
	args := []string{
		"page-content-blocks:Page Content Blocks",
		"ImageGallery:Image Gallery",
		"ImageAndTextBlock:Image And Text Block",
		"TextBlock:Text Block",
	}

	gt, err := parseType(enum.TypeFieldCollection, args)
	if err != nil {
		suite.T().Errorf("Failed: %s", err.Error())
	}

	expectedType := &item.TypeDefinition{
		Name:    "PageContentBlocks",
		Label:   "Page Content Blocks",
		Initial: "p",
		ContentBlocks: []item.ContentBlock{
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

	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedType, gt)
	}
}

func TestParseType(t *testing.T) {
	suite.Run(t, new(ParseTypeTestSuite))
}
