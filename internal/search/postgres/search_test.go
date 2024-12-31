package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/fanky5g/ponzu/database"
	pgdb "github.com/fanky5g/ponzu/internal/database/postgres"
	"github.com/fanky5g/ponzu/internal/database/postgres/connection"
	"github.com/fanky5g/ponzu/internal/database/postgres/repository"
	"github.com/fanky5g/ponzu/internal/database/postgres/test"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SearchTestSuite struct {
	suite.Suite

	client *Client
	conn   *pgxpool.Pool
}

var testData = []test.Entity{
	{Name: "Alice Smith", Email: "alice.smith@example.com", Age: 34},
	{Name: "Alicia Johnson", Email: "alicia.johnson@example.com", Age: 28},
	{Name: "Alison Brown", Email: "alison.brown@example.com", Age: 42},
	{Name: "Eve Wilson", Email: "eve.wilson@example.com", Age: 36},
	{Name: "Alfred Johnson", Email: "alfred.johnson@example.com", Age: 45},
	{Name: "Frank Johnson", Email: "frank.johnson@example.com", Age: 31},
	{Name: "David Miller", Email: "david.miller@example.com", Age: 29},
	{Name: "Alexander Johnson", Email: "alexander.johnson@example.com", Age: 40},
	{Name: "Alexa Taylor", Email: "alexa.johnson@example.com", Age: 27},
	{Name: "Alexandra Miller", Email: "alexandra.miller@example.com", Age: 33},
}

func (s *SearchTestSuite) SetupSuite() {
	var err error
	model := new(test.Model)
	db, err := pgdb.New([]database.ModelInterface{model})
	if err != nil {
		s.T().Fatal(err)
	}

	s.client, err = New(db)
	if err != nil {
		s.T().Fatal(err)
	}

	s.conn, err = connection.Get(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}

	var repo *repository.Repository
	repo, err = repository.New(s.conn, model)
	if err != nil {
		s.T().Fatal(err)
	}

	for i := range testData {
		var insert interface{}
		insert, err = repo.Insert(&testData[i])
		if err != nil {
			s.T().Fatal(err)
		}

		testData[i] = *(insert.(*test.Entity))
	}
}

func (s *SearchTestSuite) TearDownSuite() {
	ctx := context.Background()
	_, err := s.conn.Exec(ctx, "TRUNCATE TABLE test CASCADE")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *SearchTestSuite) TestSearch() {
	matches, count, err := s.client.SearchWithPagination(new(test.Entity), "Al", 0, 0)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), count, 7)
		assert.Len(s.T(), matches, 7)
	}
}

func (s *SearchTestSuite) TestSearchWithLimit() {
	matches, count, err := s.client.SearchWithPagination(new(test.Entity), "Al", 5, 0)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), 7, count)
		assert.Len(s.T(), matches, 5)
	}
}

func (s *SearchTestSuite) TestSearchWithOffset() {
	matches, count, err := s.client.SearchWithPagination(new(test.Entity), "Al", 0, 5)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), 7, count)
		assert.Len(s.T(), matches, 2)
	}
}

func TestMain(m *testing.M) {
	_ = os.Setenv("DATABASE_HOST", "localhost")
	_ = os.Setenv("DATABASE_USER", "postgres")
	_ = os.Setenv("DATABASE_PASSWORD", "password")
	_ = os.Setenv("DATABASE_NAME", "postgres")
	_ = os.Setenv("DATABASE_PORT", "5432")

	m.Run()
}

func TestSearch(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}
