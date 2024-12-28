package postgres

import (
	"context"
	"fmt"

	databasePkg "github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/internal/database/postgres/connection"
	"github.com/fanky5g/ponzu/internal/database/postgres/models"
	"github.com/fanky5g/ponzu/internal/database/postgres/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	conn         *pgxpool.Pool
	repositories map[string]database.Repository
}

func (database *Database) GetRepositoryByToken(token string) database.Repository {
	if repo, ok := database.repositories[token]; ok {
		return repo
	}

	return nil
}

func (database *Database) Close() error {
	database.conn.Close()
	return nil
}

func New(contentModels []databasePkg.ModelInterface) (*Database, error) {
	ctx := context.Background()
	conn, err := connection.Get(ctx)

	if err != nil {
		return nil, err
	}

	m := append(models.GetPonzuModels(), contentModels...)
	repos := make(map[string]database.Repository)
	for _, model := range m {
		entity := model.NewEntity()
		persistable, ok := entity.(database.Persistable)
		if !ok {
			return nil, fmt.Errorf("entity %T is not persistable", entity)
		}

		var repo database.Repository
		repo, err = repository.New(conn, model)
		if err != nil {
			return nil, err
		}

		repos[persistable.GetRepositoryToken()] = repo
	}

	d := &Database{conn: conn, repositories: repos}

	return d, nil
}
