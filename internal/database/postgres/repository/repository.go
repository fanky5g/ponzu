package repository

import (
	"context"
	"sync"

	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/database/postgres/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var (
	migrator migrations.Migrator
	once sync.Once
)

type Repository struct {
	pool     *pgxpool.Pool
	model    database.ModelInterface
}

func (repo *Repository) TableName() string {
	return repo.model.Name()
}

func New(pool *pgxpool.Pool, model database.ModelInterface) (*Repository, error) {
	var err error
	once.Do(func() {
		if migrator == nil {
			migrator, err = migrations.New(pool)
		}
	})

	if err != nil {
		return nil, err
	}

	if err = migrator.Run(context.Background(), model); err != nil {
		return nil, errors.Wrap(err, "Failed to run migration.")
	}

	repo := &Repository{
		pool:  pool,
		model: model,
	}

	return repo, nil
}
