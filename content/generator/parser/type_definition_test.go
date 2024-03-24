package parser

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/generator/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ParseTypeTestSuite struct {
	suite.Suite
	p Parser
}

func (s *ParseTypeTestSuite) SetupSuite() {
	var err error
	s.p, err = New(content.Types{
		Content:          make(map[string]content.Builder),
		FieldCollections: make(map[string]content.Builder),
		Definitions:      make(map[string]types.TypeDefinition),
	})

	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ParseTypeTestSuite) TestTypeName() {
	// blog title:string Author:string PostCategory:string entities:string some_thing:int
	args := []string{
		"blog", "title:string", "Author:string",
		"PostCategory:string", "entities:string",
		"some_thing:int", "Some_otherThing:float64",
	}

	gt, err := s.p.ParseTypeDefinition(content.TypeContent, args)
	if err != nil {
		s.T().Errorf("Failed: %s", err.Error())
	}

	if gt.Name != "Blog" {
		s.T().Errorf("Expected %s, got: %s", "Blog", gt.Name)
	}
}

func (s *ParseTypeTestSuite) TestParsing() {
	// page-content-blocks:"Page Content Blocks" ImageGallery:"Image Gallery" ImageAndTextBlock:"Image And Text Block"
	// TextBlock:"Text Block"
	args := []string{
		"page-content-blocks:Page Content Blocks",
		"ImageGallery:Image Gallery",
		"ImageAndTextBlock:Image And Text Block",
		"TextBlock:Text Block",
	}

	gt, err := s.p.ParseTypeDefinition(content.TypeFieldCollection, args)
	if err != nil {
		s.T().Errorf("Failed: %s", err.Error())
	}

	expectedType := &types.TypeDefinition{
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

	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expectedType, gt)
	}
}

func TestParseType(t *testing.T) {
	suite.Run(t, new(ParseTypeTestSuite))
}
