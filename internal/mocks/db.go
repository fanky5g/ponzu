package mocks

import (
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/stretchr/testify/mock"
)

type DB struct {
	*mock.Mock
}

func (db *DB) GetRepositoryByToken(name string) database.Repository {
	return &Repository{Mock: db.Mock}
}

func (*DB) Close() error {
	return nil
}
