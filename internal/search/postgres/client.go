package postgres

import (
	"errors"

	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/internal/database/postgres"
)

var ErrInvalidDatabase = errors.New("postgres search database invalid")

type Client struct {
	db database.Database
}

func New(db database.Database) (*Client, error) {
	_, ok := db.(*postgres.Database)
	if !ok {
		return nil, ErrInvalidDatabase
	}

	return &Client{db: db}, nil
}
