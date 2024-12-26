package repository

import (
	"errors"
	"fmt"

	"github.com/fanky5g/ponzu/database"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
)

type Row interface {
	Scan(dest ...any) error
}

type RowScanner interface {
	ScanRow(row Row) (interface{}, error)
}

func (repo *Repository) ScanRow(row Row) (interface{}, error) {
	var idBytes []byte
	var createdAt, updatedAt pgtype.Timestamptz

	model := &database.Model{
		Document: repo.model.ToDocument(repo.model.NewEntity()),
	}

	if err := row.Scan(&idBytes, &createdAt, &updatedAt, model.Document); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to scan row: %v", err)
	}

	model.CreatedAt = createdAt.Time
	model.UpdatedAt = updatedAt.Time
	model.ID = uuid.Must(uuid.FromBytes(idBytes))

	return repo.MapFromEntity(model)
}
