package editor

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ReflectionUtilsTestSuite struct {
	suite.Suite
}

func (suite *ReflectionUtilsTestSuite) TestTagNameFromStructField() {
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

func (suite *ReflectionUtilsTestSuite) TestTagNameFromStructFieldNested() {
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

func (suite *ReflectionUtilsTestSuite) TestTagNameFromStructFieldNested2() {
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

func (suite *ReflectionUtilsTestSuite) TestTagNameFromStructFieldNestedArray() {
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

func (suite *ReflectionUtilsTestSuite) TestTagNameFromStructFieldNestedArray2() {
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

func (suite *ReflectionUtilsTestSuite) TestTagNameFromStructFieldWithPositionalArg() {
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

	expectedTagName := "author.1.books.%author%.published"
	assert.Equal(
		suite.T(),
		expectedTagName,
		TagNameFromStructField(
			"Author.1.Books.%author%.Published",
			&Review{
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
			},
			&FieldArgs{PositionalPlaceHolders: []string{"%author%"}},
		),
	)
}

func (suite *ReflectionUtilsTestSuite) TestValueFromStructField() {
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

func (suite *ReflectionUtilsTestSuite) TestValueFromStructFieldNested() {
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

func (suite *ReflectionUtilsTestSuite) TestValueFromStructFieldNested2() {
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
	assert.Equal(suite.T(), []string{"tag 1", "tag 2"}, ValueFromStructField("Tags", v, nil))
	assert.Equal(suite.T(), books, ValueFromStructField("Author.Books", v, nil))
}

func TestValuesStructHelpers(t *testing.T) {
	suite.Run(t, new(ReflectionUtilsTestSuite))
}
