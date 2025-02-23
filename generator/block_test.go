package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BlockTestSuite struct {
	suite.Suite
}

func (s *BlockTestSuite) TestReferenceBlockType() {
	block := newBlock("author:@author", Field)
	assert.Equal(s.T(), "Author", block.Name)
	assert.Equal(s.T(), Field, block.Type)
	assert.Equal(s.T(), "author", block.JSONName)
	assert.Equal(s.T(), "Author", block.Label)
	assert.Equal(s.T(), "string", block.TypeName)
	assert.Equal(s.T(), "author", block.Definition.Title)
	assert.Equal(s.T(), "Author", block.ReferenceName)
	assert.Equal(s.T(), true, block.Definition.IsReference)
	assert.Equal(s.T(), false, block.Definition.IsArray)
}

func (s *BlockTestSuite) TestReferenceArrayBlockType() {
	block := newBlock("author:[]@author", Field)
	assert.Equal(s.T(), "Author", block.Name)
	assert.Equal(s.T(), Field, block.Type)
	assert.Equal(s.T(), "author", block.JSONName)
	assert.Equal(s.T(), "Author", block.Label)
	assert.Equal(s.T(), "[]string", block.TypeName)
	assert.Equal(s.T(), "author", block.Definition.Title)
	assert.Equal(s.T(), "Author", block.ReferenceName)
	assert.Equal(s.T(), true, block.Definition.IsReference)
	assert.Equal(s.T(), true, block.Definition.IsArray)
}

func (s *BlockTestSuite) TestGenerateUserDefinedBlockType() {
	tt := []struct {
		name            string
		fieldDefinition string
		expectedBlock   Block
	}{
		{
			name:            "Generate string field with type richtext",
			fieldDefinition: "description:string:richtext",
			expectedBlock: Block{
				Type:     Field,
				Name:     "Description",
				JSONName: "description",
				Label:    "Description",
				TypeName: "string",
				Definition: BlockDefinition{
					Title:       "description",
					Type:        "string:richtext",
					IsArray:     false,
					IsReference: false,
				},
			},
		},
		{
			name:            "Generate string field with select type",
			fieldDefinition: "gender:string:select@male~Male,female~Female,divers~Divers",
			expectedBlock: Block{
				Type:     Field,
				Name:     "Gender",
				JSONName: "gender",
				Label:    "Gender",
				TypeName: "string",
				Definition: BlockDefinition{
					Title:       "gender",
					Type:        "string:select",
					IsArray:     false,
					IsReference: false,
					Tokens:      []string{"male:Male", "female:Female", "divers:Divers"},
				},
			},
		},
	}

	for _, tc := range tt {
		s.T().Run(tc.name, func(t *testing.T) {
			block := newBlock(tc.fieldDefinition, Field)
			assert.Equal(s.T(), tc.expectedBlock, block)
		})
	}
}

func TestBlock(t *testing.T) {
	suite.Run(t, new(BlockTestSuite))
}
