package repository

import (
	"context"
	"fmt"
)

func (repo *Repository) Insert(entity interface{}) (interface{}, error) {
	model := repo.MapToEntity(entity)
	sqlString := fmt.Sprintf(`
INSERT INTO %s (id, created_at, updated_at, document) VALUES($1::uuid, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $2::jsonb)
`, repo.model.Name())

	ctx := context.Background()
	conn, err := repo.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	_, err = conn.Exec(
		context.Background(),
		sqlString,
		model.ID,
		model.Document,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert into %s: %v", repo.model.Name(), err)
	}

	return repo.findOneByIdWithConn(model.ID.String(), conn)
}
