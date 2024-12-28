package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (repo *Repository) GetNumberOfRows() (int, error) {
	ctx := context.Background()
	conn, err := repo.pool.Acquire(ctx)
	if err != nil {
		return 0, err
	}

	defer conn.Release()
	return repo.count(ctx, conn)
}

func (repo *Repository) count(ctx context.Context, conn *pgxpool.Conn) (int, error) {
	sqlString := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE deleted_at IS NULL
`, repo.model.Name())

	count := 0
	err := conn.QueryRow(ctx, sqlString).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
