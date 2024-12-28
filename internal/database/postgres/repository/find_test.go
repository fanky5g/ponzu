package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/internal/database/postgres/connection"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testEntities = []*testEntity{
	{
		Name:  "Foo Bar 1",
		Email: "foo@bar1.domain",
		Age:   10,
	},
	{
		Name:  "Foo Bar 2",
		Email: "foo@bar2.domain",
		Age:   20,
	},
	{
		Name:  "Foo Bar 3",
		Email: "foo@bar3.domain",
		Age:   30,
	},
}

type FindTestSuite struct {
	suite.Suite
	repo database.Repository
	conn *pgxpool.Pool
}

func (s *FindTestSuite) SetupSuite() {
	DefaultQuerySize = 3 // set to 3 to allow batching in FindAll
	conn, err := connection.Get(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}

	repo, err := New(conn, new(testModel))
	if err != nil {
		s.FailNow(err.Error())
	}

	s.repo = repo
	s.conn = conn
}

func (s *FindTestSuite) TearDownTest() {
	ctx := context.Background()
	_, err := s.conn.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", testModelToken))
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *FindTestSuite) TestFindOneById() {
	entity := &testEntity{
		Name:  "Foo Bar",
		Email: "foo@bar.domain",
		Age:   39,
	}

	ins, err := s.repo.Insert(entity)
	if err != nil {
		s.T().Fatal(err)
	}

	match, err := s.repo.FindOneById(ins.(*testEntity).ID)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), ins, match)
	}
}

func (s *FindTestSuite) TestFindOneByCriteria() {
	entity := &testEntity{
		Name:  "Foo Bar",
		Email: "foo@bar.domain",
		Age:   39,
	}

	ins, err := s.repo.Insert(entity)
	if err != nil {
		s.T().Fatal(err)
	}

	match, err := s.repo.FindOneBy(map[string]interface{}{
		"email": "foo@bar.domain",
	})
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), ins, match)
	}
}

func (s *FindTestSuite) TestFindDescNoPagination() {
	for _, entity := range testEntities {
		_, err := s.repo.Insert(entity)
		if err != nil {
			s.T().Fatal(err)
		}

		<-time.After(time.Microsecond)
	}

	numItems, matches, err := s.repo.Find(constants.Descending, 0, 0)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), len(testEntities), numItems)
		assert.Equal(s.T(), len(matches), 3)

		assert.Equal(s.T(), "foo@bar3.domain", matches[0].(*testEntity).Email)
		assert.Equal(s.T(), "foo@bar2.domain", matches[1].(*testEntity).Email)
		assert.Equal(s.T(), "foo@bar1.domain", matches[2].(*testEntity).Email)
	}
}

func (s *FindTestSuite) TestFindAscNoPagination() {
	for _, entity := range testEntities {
		_, err := s.repo.Insert(entity)
		if err != nil {
			s.T().Fatal(err)
		}

		<-time.After(time.Microsecond)
	}

	numItems, matches, err := s.repo.Find(constants.Ascending, 0, 0)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), len(testEntities), numItems)
		assert.Equal(s.T(), len(matches), 3)

		assert.Equal(s.T(), "foo@bar1.domain", matches[0].(*testEntity).Email)
		assert.Equal(s.T(), "foo@bar2.domain", matches[1].(*testEntity).Email)
		assert.Equal(s.T(), "foo@bar3.domain", matches[2].(*testEntity).Email)
	}
}

func (s *FindTestSuite) TestFindDescPagination() {
	for _, entity := range testEntities {
		_, err := s.repo.Insert(entity)
		if err != nil {
			s.T().Fatal(err)
		}

		<-time.After(time.Microsecond)
	}

	numItems, matches, err := s.repo.Find(constants.Descending, 1, 0)

	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), len(testEntities), numItems)
		assert.Equal(s.T(), len(matches), 1)
		assert.Equal(s.T(), "foo@bar3.domain", matches[0].(*testEntity).Email)
	}
}

func (s *FindTestSuite) TestFindAscPagination() {
	for _, entity := range testEntities {
		_, err := s.repo.Insert(entity)
		if err != nil {
			s.T().Fatal(err)
		}

		<-time.After(time.Microsecond)
	}

	numItems, matches, err := s.repo.Find(constants.Ascending, 1, 0)

	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), len(testEntities), numItems)
		assert.Equal(s.T(), len(matches), 1)
		assert.Equal(s.T(), "foo@bar1.domain", matches[0].(*testEntity).Email)
	}
}

func (s *FindTestSuite) TestFindAscPaginationWithOffset() {
	for _, entity := range testEntities {
		_, err := s.repo.Insert(entity)
		if err != nil {
			s.T().Fatal(err)
		}

		<-time.After(time.Microsecond)
	}

	numItems, matches, err := s.repo.Find(constants.Ascending, 2, 1)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), len(testEntities), numItems)
		assert.Equal(s.T(), len(matches), 2)
		assert.Equal(s.T(), "foo@bar2.domain", matches[0].(*testEntity).Email)
		assert.Equal(s.T(), "foo@bar3.domain", matches[1].(*testEntity).Email)
	}
}

func (s *FindTestSuite) TestFindAll() {
	for _, entity := range []*testEntity{
		testEntities[0],
		testEntities[1],
		testEntities[2],
		{
			Name:  "Foo Bar 4",
			Email: "foo@bar4.domain",
			Age:   40,
		},
		{
			Name:  "Foo Bar 5",
			Email: "foo@bar5.domain",
			Age:   50,
		},
	} {
		_, err := s.repo.Insert(entity)
		if err != nil {
			s.T().Fatal(err)
		}

		<-time.After(time.Microsecond)
	}

	matches, err := s.repo.FindAll()
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), len(matches), 5)
		assert.Equal(s.T(), "foo@bar5.domain", matches[0].(*testEntity).Email)
		assert.Equal(s.T(), "foo@bar4.domain", matches[1].(*testEntity).Email)
		assert.Equal(s.T(), "foo@bar3.domain", matches[2].(*testEntity).Email)
		assert.Equal(s.T(), "foo@bar2.domain", matches[3].(*testEntity).Email)
		assert.Equal(s.T(), "foo@bar1.domain", matches[4].(*testEntity).Email)
	}
}

func TestFind(t *testing.T) {
	suite.Run(t, new(FindTestSuite))
}
