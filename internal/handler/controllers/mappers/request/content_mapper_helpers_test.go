package request

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
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

func TestContentMapperHelpers(t *testing.T) {
	suite.Run(t, new(ContentMapperHelpersTestSuite))
}
