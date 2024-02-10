package request

import (
	"github.com/fanky5g/ponzu/internal/domain/entities"
	"github.com/fanky5g/ponzu/internal/domain/enum"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/url"
	"testing"
)

type MapSearchRequestTestSuite struct {
	suite.Suite
}

func (suite *MapSearchRequestTestSuite) TestGetSearchRequestDtoWithEmptyQueryValues() {
	var q url.Values

	expectedSearchRequestDto := &request.SearchRequestDto{
		SortOrder: enum.Descending,
		PaginationRequestDto: request.PaginationRequestDto{
			Count:  10,
			Offset: 0,
		},
	}

	searchRequestDto, err := GetSearchRequestDto(q)
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

	expectedSearchRequestDto := &request.SearchRequestDto{
		Query:     "Alpaka",
		SortOrder: enum.Ascending,
		PaginationRequestDto: request.PaginationRequestDto{
			Count:  100,
			Offset: 5,
		},
	}

	searchRequestDto, err := GetSearchRequestDto(q)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedSearchRequestDto, searchRequestDto)
	}
}

func (suite *MapSearchRequestTestSuite) TestMapSearchRequest() {
	searchRequestDto := &request.SearchRequestDto{
		Query:     "Alpaka",
		SortOrder: enum.Ascending,
		PaginationRequestDto: request.PaginationRequestDto{
			Count:  100,
			Offset: 5,
		},
	}

	expectedSearch := &entities.Search{
		Query:     "Alpaka",
		SortOrder: enum.Ascending,
		Pagination: &entities.Pagination{
			Count:  100,
			Offset: 5,
		},
	}

	search, err := MapSearchRequest(searchRequestDto)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedSearch, search)
	}
}

func TestMapSearchRequest(t *testing.T) {
	suite.Run(t, new(MapSearchRequestTestSuite))
}
