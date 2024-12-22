package request

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilTestSuite struct {
	suite.Suite
}

func (suite *UtilTestSuite) TestMapJSONContentToURLValues() {
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

func (suite *UtilTestSuite) TestMapNestedJSONContentToURLValues() {
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
