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
		TagNameFromStructField("Title", &Review{}, nil),
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
		TagNameFromStructField("Author.Name", &Review{}, nil),
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
		TagNameFromStructField("Author.Books.Published", &Review{}, nil),
	)
}

func (suite *ValuesTestSuite) TestTagNameFromStructFieldNestedArray() {
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

	expectedTagName := "author.books.0.published"
	assert.Equal(
		suite.T(),
		expectedTagName,
		TagNameFromStructField("Author.Books.0.Published", &Review{
			Title: "Title",
			Author: Author{
				Name: "Author",
				Books: []Book{
					{
						Name:      "Book",
						Published: "Date",
					},
				},
			},
		}, nil),
	)
}

func (suite *ValuesTestSuite) TestTagNameFromStructFieldNestedArray2() {
	type Book struct {
		Name      string `json:"name"`
		Published string `json:"published"`
	}

	type Author struct {
		Name  string `json:"name"`
		Books []Book `json:"books"`
	}

	type Review struct {
		Title  string   `json:"title"`
		Author []Author `json:"author"`
	}

	expectedTagName := "author.1.books.0.published"
	assert.Equal(
		suite.T(),
		expectedTagName,
		TagNameFromStructField("Author.1.Books.0.Published", &Review{
			Title: "Title",
			Author: []Author{
				{
					Name: "Author",
					Books: []Book{
						{
							Name:      "Book",
							Published: "Date",
						},
					},
				},
			},
		}, nil),
	)
}

func (suite *ValuesTestSuite) TestTagNameFromStructFieldWithPositionalArg() {
	type Book struct {
		Name      string `json:"name"`
		Published string `json:"published"`
	}

	type Author struct {
		Name  string `json:"name"`
		Books []Book `json:"books"`
	}

	type Review struct {
		Title  string   `json:"title"`
		Author []Author `json:"author"`
	}

	expectedTagName := "author.1.books.%pos%.published"
	assert.Equal(
		suite.T(),
		expectedTagName,
		TagNameFromStructField("Author.1.Books.%pos%.Published", &Review{
			Title: "Title",
			Author: []Author{
				{
					Name: "Author",
					Books: []Book{
						{
							Name:      "Book",
							Published: "Date",
						},
					},
				},
			},
		}, nil),
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
		ValueFromStructField("Title", Review{Title: title}, nil),
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
		ValueFromStructField(
			"Author.Name",
			&Review{Title: "Review Title", Author: Author{Name: name}},
			nil,
		),
	)
}

func (suite *ValuesTestSuite) TestValueFromStructFieldNested2() {
	type Book struct {
		Name      string `json:"name"`
		Published string `json:"published"`
	}

	type Author struct {
		Name  string `json:"name"`
		Books []Book `json:"books"`
	}

	type Review struct {
		Title  string   `json:"title"`
		Author Author   `json:"author"`
		Tags   []string `json:"tags"`
	}

	publishedDate := "BookPublishedDate"
	books := []Book{
		{
			Name:      "BookName",
			Published: publishedDate,
		},
	}

	v := &Review{
		Title: "ReviewTitle",
		Author: Author{
			Name:  "AuthorName",
			Books: books,
		},
		Tags: []string{"tag 1", "tag 2"},
	}

	assert.Equal(suite.T(), publishedDate, ValueFromStructField("Author.Books.0.Published", v, nil))
	// legacy ponzu string array joins
	assert.Equal(suite.T(), "tag 1__ponzutag 2", ValueFromStructField("Tags", v, nil))
	assert.Equal(suite.T(), books, ValueFromStructField("Author.Books", v, nil))
}

func TestValuesStructHelpers(t *testing.T) {
	suite.Run(t, new(ValuesTestSuite))
}
