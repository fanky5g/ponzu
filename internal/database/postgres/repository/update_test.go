package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/fanky5g/ponzu/internal/database/postgres/connection"
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UpdateTestSuite struct {
	suite.Suite
	repo   database.Repository
	conn   *pgxpool.Pool
	entity *testEntity
}

func (s *UpdateTestSuite) SetupSuite() {
	conn, err := connection.Get(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}

	repo, err := New(conn, new(testModel))
	if err != nil {
		s.FailNow(err.Error())
	}

	entity, err := repo.Insert(&testEntity{
		ID:    uuid.New().String(),
		Name:  "Foo Bar",
		Email: "foo@bar.domain",
		Age:   39,
	})

	if err != nil {
		s.T().Fatal(err)
	}

	s.entity = entity.(*testEntity)
	s.repo = repo
	s.conn = conn
}

func (s *UpdateTestSuite) TearDownSuite() {
	ctx := context.Background()
	_, err := s.conn.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", testModelToken))
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *UpdateTestSuite) TestUpdateNonExistingEntity() {
	nonExistingEntityId := uuid.New().String()
	update := &testEntity{
		ID:    nonExistingEntityId,
		Name:  "Foo Baz",
		Email: "foo@baz.domain",
		Age:   20,
	}

	updated, err := s.repo.UpdateById(nonExistingEntityId, update)
	assert.Nil(s.T(), updated)
	assert.Error(s.T(), err)
}

func (s *UpdateTestSuite) TestUpdateExistingEntity() {
	existingEntityId := s.entity.ID
	update := &testEntity{
		ID:    existingEntityId,
		Name:  "Foo Baz",
		Email: "foo@baz.domain",
		Age:   20,
	}

	u, err := s.repo.UpdateById(existingEntityId, update)
	if assert.NoError(s.T(), err) {
		updated := u.(*testEntity)
		assert.Equal(s.T(), updated.ID, existingEntityId)
		assert.Equal(s.T(), updated.Name, update.Name)
		assert.Equal(s.T(), updated.Email, update.Email)
		assert.Equal(s.T(), updated.Age, update.Age)
	}
}

func TestUpdate(t *testing.T) {
	suite.Run(t, new(UpdateTestSuite))
}
