package editor

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ValuesTestSuite struct {
	suite.Suite
}

func (suite *ValuesTestSuite) TestTagNameFromStructField() {
	type Review struct {
		Title string `json:"title"`
	}

	expectedTagName := "title"
	assert.Equal(
		suite.T(),
		expectedTagName,
		TagNameFromStructField("Title", &Review{}),
	)
}

func (suite *ValuesTestSuite) TestTagNameFromStructFieldNested() {
	type Author struct {
		Name string `json:"name"`
	}

	type Review struct {
		Title  string `json:"title"`
		Author Author `json:"author"`
	}

	expectedTagName := "author.name"
	assert.Equal(
		suite.T(),
		expectedTagName,
		TagNameFromStructField("Author.Name", &Review{}),
	)
}

func (suite *ValuesTestSuite) TestTagNameFromStructFieldNested2() {
	type Book struct {
		Name      string `json:"name"`
		Published string `json:"published"`
	}

	type Author struct {
		Name  string `json:"name"`
		Books []Book `json:"books"`
	}

	type Review struct {
		Title  string `json:"title"`
		Author Author `json:"author"`
	}

	expectedTagName := "author.books.published"
	assert.Equal(
		suite.T(),
		expectedTagName,
		TagNameFromStructField("Author.Books.Published", &Review{}),
	)
}

func (suite *ValuesTestSuite) TestValueFromStructField() {
	type Review struct {
		Title string `json:"title"`
	}

	title := "Review Title"
	assert.Equal(
		suite.T(),
		title,
		ValueFromStructField("Title", Review{Title: title}),
	)
}

func (suite *ValuesTestSuite) TestValueFromStructFieldNested() {
	type Author struct {
		Name string `json:"name"`
	}

	type Review struct {
		Title  string `json:"title"`
		Author Author `json:"author"`
	}

	name := "Ponzu"
	assert.Equal(
		suite.T(),
		name,
		ValueFromStructField("Author.Name", &Review{Title: "Review Title", Author: Author{Name: name}}),
	)
}

func TestValuesStructHelpers(t *testing.T) {
	suite.Run(t, new(ValuesTestSuite))
}
