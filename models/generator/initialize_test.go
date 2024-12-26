package generator

import (
	"github.com/fanky5g/ponzu/generator"
	"github.com/stretchr/testify/assert"
	"go/format"
)

func (s *GeneratorTestSuite) TestInitialize() {
	t := &generator.TypeDefinition{
		Type:  generator.Content,
		Name:  "Author",
		Label: "Author",
	}

	w := new(testWriter)

	expectedBuffer, err := format.Source([]byte(`
// Code generated by ponzu. DO NOT EDIT.
package models

import "github.com/fanky5g/ponzu/database"

var Models = make([]database.ModelInterface, 0)
	`))

	if err != nil {
		s.T().Fatal(err)
	}

	if assert.NoError(s.T(), s.g.Initialize(t, w)) {
		assert.Equal(s.T(), expectedBuffer, w.buf)
	}
}
