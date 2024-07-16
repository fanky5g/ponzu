package infrastructure

import (
	"fmt"
	"os"
	"testing"

	bleveSearch "github.com/fanky5g/ponzu-driver-bleve"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SearchClientTestSuite struct {
	suite.Suite
}

func (s *SearchClientTestSuite) TestGetBleveSearchClient() {
	tempDir, err := os.MkdirTemp(os.TempDir(), "bleve")
	if err != nil {
		s.T().Fatal(err)
		return
	}

	os.Args = append(os.Args, "--search_driver=bleve")
	os.Args = append(os.Args, fmt.Sprintf("--data_dir=%s", tempDir))

	searchClient, err := getSearchClient()

	if assert.NoError(s.T(), err) {
		assert.IsType(s.T(), &bleveSearch.Client{}, searchClient)
	}
}

func TestGetSearchClient(t *testing.T) {
	suite.Run(t, new(SearchClientTestSuite))
}
