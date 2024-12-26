package models

import (
	"github.com/fanky5g/ponzu/generator"
	"github.com/stretchr/testify/assert"
	"go/format"
)

func (s *GeneratorTestSuite) TestGenerate() {
	t := &generator.TypeDefinition{
		Type:  generator.Content,
		Name:  "Author",
		Label: "Author",
	}

	w := new(testWriter)

	expectedBuffer, err := format.Source([]byte(`
package models

import (
	"encoding/json"
	"fmt"
	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/testapp/entities"
	"strings"
)

type AuthorDocument struct {
	*entities.Author
}

func (document *AuthorDocument) Value() (interface{}, error) {
	return json.Marshal(document)
}

func (document *AuthorDocument) Scan(src interface{}) error {
	if byteSrc, ok := src.([]byte); ok {
		return json.Unmarshal(byteSrc, &document)
	}

	if stringSrc, ok := src.(string); ok {
		return json.NewDecoder(strings.NewReader(stringSrc)).Decode(&document)
	}

	return fmt.Errorf("unsupported type %T", src)
}

type AuthorModel struct{}

func (*AuthorModel) Name() string {
	return "author"
}

func (*AuthorModel) NewEntity() interface{} {
	return new(entities.Author)
}

func (model *AuthorModel) ToDocument(entity interface{}) database.DocumentInterface {
	return &AuthorDocument{
		Author: entity.(*entities.Author),
	}
}

func init() {
	Models = append(Models, new(AuthorModel))
}
	`))

	if err != nil {
		s.T().Fatal(err)
	}

	if assert.NoError(s.T(), s.g.Generate(t, w)) {
		assert.Equal(s.T(), expectedBuffer, w.buf)
	}
}
