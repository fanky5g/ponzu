package request

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/test/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/url"
	"testing"
)

type page struct {
	item.Item

	Title         string             `json:"title"`
	URL           string             `json:"url"`
	ContentBlocks *pageContentBlocks `json:"content_blocks"`
}

type pageContentBlocks []content.FieldCollection

type textBlock struct {
	Text string `json:"text"`
}

func (p *pageContentBlocks) Name() string {
	return "Page Content Blocks"
}

func (p *pageContentBlocks) Data() []content.FieldCollection {
	return *p
}

func (p *pageContentBlocks) Add(fieldCollection content.FieldCollection) {
	*p = append(*p, fieldCollection)
}

func (p *pageContentBlocks) Set(i int, fieldCollection content.FieldCollection) {
	data := p.Data()
	data[i] = fieldCollection
	*p = data
}

func (p *pageContentBlocks) SetData(data []content.FieldCollection) {
	*p = data
}

func (p *pageContentBlocks) AllowedTypes() map[string]content.Builder {
	return map[string]content.Builder{
		"TextBlock": func() interface{} {
			return new(textBlock)
		},
	}
}

type ContentMapperHelpersTestSuite struct {
	suite.Suite
}

func (suite *ContentMapperHelpersTestSuite) TestMapPayloadToGenericEntity() {
	type Review struct {
		item.Item

		Title  string   `json:"title"`
		Body   string   `json:"body"`
		Rating int      `json:"rating"`
		Tags   []string `json:"tags"`
	}

	payload := url.Values{
		"id":        []string{"6"},
		"uuid":      []string{"183a4535-f015-4660-bb8f-6541522e9afb"},
		"body":      []string{"API Content Body"},
		"rating":    []string{"20"},
		"slug":      []string{"review-183a4535-f015-4660-bb8f-6541522e9afb"},
		"tags.0":    []string{"API"},
		"tags.1":    []string{"Ponzu"},
		"timestamp": []string{"1707647434000"},
		"updated":   []string{"1707647434000"},
		"title":     []string{"API Content Title"},
	}

	expectedEntity := &Review{
		Item: item.Item{
			ID:        "6",
			Slug:      "review-183a4535-f015-4660-bb8f-6541522e9afb",
			Timestamp: 1707647434000,
			Updated:   1707647434000,
		},
		Title:  "API Content Title",
		Body:   "API Content Body",
		Rating: 20,
		Tags:   []string{"API", "Ponzu"},
	}

	entity, err := MapPayloadToGenericEntity(&Review{}, payload)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedEntity, entity)
	}
}

func (suite *ContentMapperHelpersTestSuite) TestMapPayloadToGenericEntityNestedStruct() {
	type Author struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type Review struct {
		item.Item

		Title  string   `json:"title"`
		Body   string   `json:"body"`
		Rating int      `json:"rating"`
		Tags   []string `json:"tags"`
		Author Author   `json:"author"`
	}

	payload := url.Values{
		"id":          []string{"6"},
		"uuid":        []string{"183a4535-f015-4660-bb8f-6541522e9afb"},
		"body":        []string{"API Content Body"},
		"rating":      []string{"20"},
		"slug":        []string{"review-183a4535-f015-4660-bb8f-6541522e9afb"},
		"tags.0":      []string{"API"},
		"tags.1":      []string{"Ponzu"},
		"timestamp":   []string{"1707647434000"},
		"updated":     []string{"1707647434000"},
		"title":       []string{"API Content Title"},
		"author.age":  []string{"25"},
		"author.name": []string{"Foo Bar"},
	}

	expectedEntity := &Review{
		Item: item.Item{
			ID:        "6",
			Slug:      "review-183a4535-f015-4660-bb8f-6541522e9afb",
			Timestamp: 1707647434000,
			Updated:   1707647434000,
		},
		Title:  "API Content Title",
		Body:   "API Content Body",
		Rating: 20,
		Tags:   []string{"API", "Ponzu"},
		Author: Author{
			Name: "Foo Bar",
			Age:  25,
		},
	}

	entity, err := MapPayloadToGenericEntity(&Review{}, payload)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedEntity, entity)
	}
}

func (suite *ContentMapperHelpersTestSuite) TestMapPayloadToGenericEntityNestedStruct2() {
	type Author struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type Review struct {
		item.Item

		Title  string   `json:"title"`
		Body   string   `json:"body"`
		Rating int      `json:"rating"`
		Tags   []string `json:"tags"`
		Author []Author `json:"authors"`
	}

	payload := url.Values{
		"id":                              []string{"6"},
		"uuid":                            []string{"183a4535-f015-4660-bb8f-6541522e9afb"},
		"body":                            []string{"API Content Body"},
		"rating":                          []string{"20"},
		"slug":                            []string{"review-183a4535-f015-4660-bb8f-6541522e9afb"},
		"tags.0":                          []string{"API"},
		"tags.1":                          []string{"Ponzu"},
		"timestamp":                       []string{"1707647434000"},
		"updated":                         []string{"1707647434000"},
		"title":                           []string{"API Content Title"},
		"authors.0.age":                   []string{"25"},
		"authors.0.name":                  []string{"Foo Bar"},
		"authors.3.age":                   []string{"30"},
		"authors.3.name":                  []string{"Foo Bar 3"},
		"authors.5.age":                   []string{"50"},
		"authors.5.name":                  []string{"Foo Bar 5"},
		".__ponzu-repeat.authors.length":  []string{"3"},
		".__ponzu-repeat.authors.removed": []string{"1,2,4"},
	}

	expectedEntity := &Review{
		Item: item.Item{
			ID:        "6",
			Slug:      "review-183a4535-f015-4660-bb8f-6541522e9afb",
			Timestamp: 1707647434000,
			Updated:   1707647434000,
		},
		Title:  "API Content Title",
		Body:   "API Content Body",
		Rating: 20,
		Tags:   []string{"API", "Ponzu"},
		Author: []Author{
			{
				Name: "Foo Bar",
				Age:  25,
			},
			{
				Name: "Foo Bar 3",
				Age:  30,
			},
			{
				Name: "Foo Bar 5",
				Age:  50,
			},
		},
	}

	entity, err := MapPayloadToGenericEntity(&Review{}, payload)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedEntity, entity)
	}
}

func (suite *ContentMapperHelpersTestSuite) TestMapPayloadToGenericEntityNestedStruct3() {
	type Product struct {
		item.Item

		Name        string   `json:"name"`
		Description string   `json:"description"`
		Images      []string `json:"images" reference:"Photo"`
	}

	payload := url.Values{
		".__ponzu-repeat.images.length":  []string{"3"},
		".__ponzu-repeat.images.removed": []string{"1"},
		"description":                    []string{"<p>Product Description</p>"},
		"id":                             []string{"f62df6d9-9031-44a4-8db2-f73a5d5977db"},
		"images-selected":                []string{"60cbe180-3610-47de-a1eb-ac2cdb6bf4c5"},
		"images.0":                       []string{"1ee2a499-12cc-4bd0-9394-5f918b7785cd"},
		"images.2":                       []string{"35144bc4-e543-4423-9f18-e5d6f0d38176"},
		"images.3":                       []string{"60cbe180-3610-47de-a1eb-ac2cdb6bf4c5"},
		"name":                           []string{"Product Name"},
		"timestamp":                      []string{"1738212589056"},
		"updated":                        []string{"1738212589056"},
	}

	expectedEntity := &Product{
		Item: item.Item{
			ID:        "f62df6d9-9031-44a4-8db2-f73a5d5977db",
			Timestamp: 1738212589056,
			Updated:   1738212589056,
		},
		Name:        "Product Name",
		Description: "<p>Product Description</p>",
		Images: []string{
			"1ee2a499-12cc-4bd0-9394-5f918b7785cd",
			"35144bc4-e543-4423-9f18-e5d6f0d38176",
			"60cbe180-3610-47de-a1eb-ac2cdb6bf4c5",
		},
	}

	entity, err := MapPayloadToGenericEntity(&Product{}, payload)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedEntity, entity)
	}
}

func (suite *ContentMapperHelpersTestSuite) TestMapPayloadToGenericEntityFieldCollections() {
	payload := url.Values{
		"id":                          []string{"6"},
		"uuid":                        []string{"183a4535-f015-4660-bb8f-6541522e9afb"},
		"slug":                        []string{"page-183a4535-f015-4660-bb8f-6541522e9afb"},
		"timestamp":                   []string{"1707647434000"},
		"updated":                     []string{"1707647434000"},
		"title":                       []string{"Home"},
		"url":                         []string{"https://ponzu.domain"},
		"content_blocks.0.type":       []string{"TextBlock"},
		"content_blocks.0.value.text": []string{"This is some WYSIWYG entities"},
	}

	expectedEntity := &page{
		Item: item.Item{
			ID:        "6",
			Slug:      "page-183a4535-f015-4660-bb8f-6541522e9afb",
			Timestamp: 1707647434000,
			Updated:   1707647434000,
		},
		Title: "Home",
		URL:   "https://ponzu.domain",
		ContentBlocks: &pageContentBlocks{
			{
				Type:  "TextBlock",
				Value: &textBlock{Text: "This is some WYSIWYG entities"},
			},
		},
	}

	entity, err := MapPayloadToGenericEntity(&page{}, payload)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedEntity, entity)
	}
}

func (suite *ContentMapperHelpersTestSuite) TestMapPayloadToGenericEntityFieldCollections2() {
	payload := url.Values{
		"id":                          []string{"6"},
		"uuid":                        []string{"183a4535-f015-4660-bb8f-6541522e9afb"},
		"slug":                        []string{"page-183a4535-f015-4660-bb8f-6541522e9afb"},
		"timestamp":                   []string{"1707647434000"},
		"updated":                     []string{"1707647434000"},
		"title":                       []string{"Home"},
		"url":                         []string{"https://ponzu.domain"},
		"content_blocks.0.type":       []string{"TextBlock"},
		"content_blocks.0.value.text": []string{"This is some WYSIWYG entities"},
		"content_blocks.3.type":       []string{"TextBlock"},
		"content_blocks.3.value.text": []string{"This is some WYSIWYG entities 3"},
		"content_blocks.5.type":       []string{"TextBlock"},
		"content_blocks.5.value.text": []string{"This is some WYSIWYG entities 5"},
		".__ponzu-field-collection.content_blocks.length":  []string{"3"},
		".__ponzu-field-collection.content_blocks.removed": []string{"1,2,4"},
	}

	expectedEntity := &page{
		Item: item.Item{
			ID:        "6",
			Slug:      "page-183a4535-f015-4660-bb8f-6541522e9afb",
			Timestamp: 1707647434000,
			Updated:   1707647434000,
		},
		Title: "Home",
		URL:   "https://ponzu.domain",
		ContentBlocks: &pageContentBlocks{
			{
				Type:  "TextBlock",
				Value: &textBlock{Text: "This is some WYSIWYG entities"},
			},
			{
				Type:  "TextBlock",
				Value: &textBlock{Text: "This is some WYSIWYG entities 3"},
			},
			{
				Type:  "TextBlock",
				Value: &textBlock{Text: "This is some WYSIWYG entities 5"},
			},
		},
	}

	entity, err := MapPayloadToGenericEntity(&page{}, payload)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedEntity, entity)
	}
}

func TestMapPayloadToGenericEntity(t *testing.T) {
	tt := []struct {
		name         string
		Payload      url.Values
		ExpectedPage *types.Page
	}{
		{
			name: "DeleteContentBlocksEntryInPosition",
			Payload: url.Values{
				"id":        {"7219db7f-6fa5-470b-bc82-911fc6ac70c2"},
				"title":     {"HomePage"},
				"url":       {"https://ponzu.domain"},
				"slug":      {"homepage"},
				"type":      {"Page"},
				"timestamp": {"1707647434000"},
				"updated":   {"1707647434000"},
				".__ponzu-repeat.content_blocks.0.value.cta.length":  {"1"},
				".__ponzu-repeat.content_blocks.0.value.cta.removed": {"0"},
				"content_blocks.0.type":                              {"Banner"},
				"content_blocks.0.value.background.alt":              {"Men's Suit Watch"},
				"content_blocks.0.value.background.file":             {"8f1ce16b-0736-4153-b79f-971ab80997e8"},
				"content_blocks.0.value.background.position":         {"right"},
				"content_blocks.0.value.cta.1.link.href":             {"https://link-to-another-button"},
				"content_blocks.0.value.cta.1.link.type":             {"external"},
				"content_blocks.0.value.cta.1.text":                  {"Another Button"},
				"content_blocks.0.value.cta.1.type":                  {"text"},
				"content_blocks.0.value.text":                        {"Banner Text"},
			},
			ExpectedPage: &types.Page{
				Item: item.Item{
					ID:        "7219db7f-6fa5-470b-bc82-911fc6ac70c2",
					Slug:      "homepage",
					Timestamp: 1707647434000,
					Updated:   1707647434000,
				},
				Title: "HomePage",
				URL:   "https://ponzu.domain",
				ContentBlocks: &types.PageContentBlocks{
					{
						Type: "Banner",
						Value: &types.Banner{
							Background: types.BackgroundImage{
								File:     "8f1ce16b-0736-4153-b79f-971ab80997e8",
								Alt:      "Men's Suit Watch",
								Position: "right",
							},
							Text: "Banner Text",
							Cta: []types.ButtonLink{
								{
									Type: "text",
									Text: "Another Button",
									Link: types.LinkWithType{
										Href: "https://link-to-another-button",
										Type: "external",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			entity, err := MapPayloadToGenericEntity(&types.Page{}, tc.Payload)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.ExpectedPage, entity)
			}
		})
	}
}

func TestContentMapperHelpers(t *testing.T) {
	suite.Run(t, new(ContentMapperHelpersTestSuite))
}
