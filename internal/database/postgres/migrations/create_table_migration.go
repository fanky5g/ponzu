package migrations

import (
	"context"
	"fmt"

	"github.com/fanky5g/ponzu/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type createTable struct{}

func (m *createTable) Up(ctx context.Context, conn *pgxpool.Conn, model database.ModelInterface) error {
	sqlString := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	id UUID PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	deleted_at TIMESTAMPTZ NULL,
	document jsonb NOT NULL
);
`, model.Name())

	_, err := conn.Exec(ctx, sqlString)
	return err
}

func (m *createTable) Down(ctx context.Context, conn *pgxpool.Conn, model database.ModelInterface) error {
	_, err := conn.Exec(
		ctx,
		fmt.Sprintf("DROP TABLE IF EXISTS %s;", model.Name()),
	)

	return err
}
