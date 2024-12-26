package repository

import (
	"context"
	"fmt"
)

func (repo *Repository) UpdateById(id string, update interface{}) (interface{}, error) {
	document := repo.model.ToDocument(update)
	sqlString := fmt.Sprintf(
		"UPDATE %s SET document = $1::jsonb, updated_at = CURRENT_TIMESTAMP WHERE id = $2::uuid",
		repo.model.Name(),
	)

	ctx := context.Background()
	conn, err := repo.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	ct, err := conn.Exec(
		ctx,
		sqlString,
		document,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update %s[%s]: %v", repo.model.Name(), id, err)
	}

	if ct.RowsAffected() == 0 {
		return nil, fmt.Errorf("entity with id %s not found", id)
	}

	return repo.findOneByIdWithConn(id, conn)
}
