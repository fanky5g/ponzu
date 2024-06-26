package request

import (
	"bytes"
	"encoding/json"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
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

func (suite *ContentMapperHelpersTestSuite) TestMapJSONContentToURLValues() {
	request := map[string]interface{}{
		"title":     "API Content Title",
		"body":      "API Content Value",
		"rating":    20,
		"tags":      []string{"API", "Ponzu"},
		"trustable": true,
	}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(request); err != nil {
		suite.FailNow(err.Error())
	}

	req, _ := http.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")

	expectedContent := map[string][]string{
		"title":     {"API Content Title"},
		"body":      {"API Content Value"},
		"rating":    {"20"},
		"tags":      {"API", "Ponzu"},
		"trustable": {"true"},
	}

	jsonContent, err := mapJSONContentToURLValues(req)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedContent, jsonContent)
	}
}

func (suite *ContentMapperHelpersTestSuite) TestMapNestedJSONContentToURLValues() {
	request := map[string]interface{}{
		"title":     "API Content Title",
		"body":      "API Content Value",
		"rating":    20,
		"tags":      []string{"API", "Ponzu"},
		"trustable": true,
		"author": map[string]interface{}{
			"name": "Ponzu",
			"age":  25,
			"location": map[string]interface{}{
				"country":  "USA",
				"timezone": "PST",
			},
		},
	}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(request); err != nil {
		suite.FailNow(err.Error())
	}

	req, _ := http.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")

	expectedContent := map[string][]string{
		"title":                    {"API Content Title"},
		"body":                     {"API Content Value"},
		"rating":                   {"20"},
		"tags":                     {"API", "Ponzu"},
		"trustable":                {"true"},
		"author.name":              {"Ponzu"},
		"author.age":               {"25"},
		"author.location.country":  {"USA"},
		"author.location.timezone": {"PST"},
	}

	mapped, err := mapJSONContentToURLValues(req)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedContent, mapped)
	}
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

	var t content.Builder = func() interface{} {
		return new(Review)
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

	entity, err := mapPayloadToGenericEntity(t, payload)
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

	var t content.Builder = func() interface{} {
		return new(Review)
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

	entity, err := mapPayloadToGenericEntity(t, payload)
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

	var t content.Builder = func() interface{} {
		return new(Review)
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

	entity, err := mapPayloadToGenericEntity(t, payload)
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

	var t content.Builder = func() interface{} {
		return new(page)
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

	entity, err := mapPayloadToGenericEntity(t, payload)
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

	var t content.Builder = func() interface{} {
		return new(page)
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

	entity, err := mapPayloadToGenericEntity(t, payload)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedEntity, entity)
	}
}

func TestContentMapperHelpers(t *testing.T) {
	suite.Run(t, new(ContentMapperHelpersTestSuite))
}
