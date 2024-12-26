package postgres

import (
	"github.com/fanky5g/ponzu/internal/database/postgres/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SearchableFieldsTestSuite struct {
	suite.Suite
}

func (s *SearchableFieldsTestSuite) TestGetSearchableAttributes() {
	searchableFields, err := getSearchableFields(new(test.Entity))
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), searchableFields, []string{"name", "email"})
	}
}

func TestSearchableFields(t *testing.T) {
	suite.Run(t, new(SearchableFieldsTestSuite))
}
