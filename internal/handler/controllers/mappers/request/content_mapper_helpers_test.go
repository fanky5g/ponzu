package request

import (
	"bytes"
	"encoding/json"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"testing"
)

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

	content, err := mapJSONContentToURLValues(req)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedContent, content)
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

	content, err := mapJSONContentToURLValues(req)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedContent, content)
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

	var t item.EntityBuilder = func() interface{} {
		return new(Review)
	}

	uid, err := uuid.FromString("183a4535-f015-4660-bb8f-6541522e9afb")
	if err != nil {
		suite.FailNow(err.Error())
		return
	}

	expectedEntity := &Review{
		Item: item.Item{
			ID:        "6",
			UUID:      uid,
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

	//type Author struct {
	//	Name string `json:"name"`
	//	Age  int    `json:"age"`
	//}
	//"author.age": []string{"25"},
	//	"author.name": []string{"Foo Bar"},

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

	var t item.EntityBuilder = func() interface{} {
		return new(Review)
	}

	uid, err := uuid.FromString("183a4535-f015-4660-bb8f-6541522e9afb")
	if err != nil {
		suite.FailNow(err.Error())
		return
	}

	expectedEntity := &Review{
		Item: item.Item{
			ID:        "6",
			UUID:      uid,
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

func TestContentMapperHelpers(t *testing.T) {
	suite.Run(t, new(ContentMapperHelpersTestSuite))
}
