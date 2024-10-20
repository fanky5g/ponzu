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

func TestBlock(t *testing.T) {
	suite.Run(t, new(BlockTestSuite))
}
