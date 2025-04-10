package search

import (
	"bytes"
	"encoding/json"
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"testing"
)

type MapSearchRequestTestSuite struct {
	suite.Suite
}

func (suite *MapSearchRequestTestSuite) TestGetSearchRequestDtoWithEmptyQueryValues() {
	var q url.Values

	expectedSearchRequestDto := &RequestDto{
		SortOrder: constants.Descending,
		PaginationRequestDto: PaginationRequestDto{
			Count:  RowsPerPage,
			Offset: 0,
		},
	}

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "services/x-www-form-urlencoded")
	req.URL.RawQuery = q.Encode()

	searchRequestDto, err := GetSearchRequestDto(req)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedSearchRequestDto, searchRequestDto)
	}
}

func (suite *MapSearchRequestTestSuite) TestGetSearchRequestDto() {
	q := make(url.Values)

	q.Set("q", "Alpaka")
	q.Set("count", "100")
	q.Set("offset", "5")
	q.Set("order", "asc")

	expectedSearchRequestDto := &RequestDto{
		Query:     "Alpaka",
		SortOrder: constants.Ascending,
		PaginationRequestDto: PaginationRequestDto{
			Count:  100,
			Offset: 5,
		},
	}

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "services/x-www-form-urlencoded")
	req.URL.RawQuery = q.Encode()

	searchRequestDto, err := GetSearchRequestDto(req)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedSearchRequestDto, searchRequestDto)
	}
}

func (suite *MapSearchRequestTestSuite) TestGetSearchRequestDtoFromJSONRequest() {
	payload := &RequestDto{
		Query:     "Alpaka",
		SortOrder: "asc",
		PaginationRequestDto: PaginationRequestDto{
			Offset: 5,
			Count:  100,
		},
	}

	expectedSearchRequestDto := &RequestDto{
		Query:     "Alpaka",
		SortOrder: constants.Ascending,
		PaginationRequestDto: PaginationRequestDto{
			Count:  100,
			Offset: 5,
		},
	}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		assert.FailNow(suite.T(), err.Error())
		return
	}

	req, _ := http.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")

	searchRequestDto, err := GetSearchRequestDto(req)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedSearchRequestDto, searchRequestDto)
	}
}

func (suite *MapSearchRequestTestSuite) TestMapSearchRequest() {
	searchRequestDto := &RequestDto{
		Query:     "Alpaka",
		SortOrder: constants.Ascending,
		PaginationRequestDto: PaginationRequestDto{
			Count:  100,
			Offset: 5,
		},
	}

	expectedSearch := &Search{
		Query:     "Alpaka",
		SortOrder: constants.Ascending,
		Count:     100,
		Offset:    5,
	}

	search, err := MapSearchRequest(searchRequestDto)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedSearch, search)
	}
}

func TestMapSearchRequest(t *testing.T) {
	suite.Run(t, new(MapSearchRequestTestSuite))
}
