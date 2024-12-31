package migrations

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/fanky5g/ponzu/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type migrator struct {
	pool *pgxpool.Pool
}

type Migrator interface {
	Run(ctx context.Context, model database.ModelInterface) error
}

type Entry interface {
	Up(context.Context, *pgxpool.Conn, database.ModelInterface) error
	Down(context.Context, *pgxpool.Conn, database.ModelInterface) error
}

// Supported time layout: RFC3339
var migrations = map[string]Entry{
	"2024-05-20T12:00:00+02:00": new(createTable),
}

func New(pool *pgxpool.Pool) (Migrator, error) {
	if err := createMigrationTable(pool); err != nil {
		return nil, err
	}

	return &migrator{pool: pool}, nil
}

func createMigrationTable(pool *pgxpool.Pool) error {
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	sqlString := `
CREATE TABLE IF NOT EXISTS migrations (
	id varchar(256) PRIMARY KEY,
	executed_at TIMESTAMP NOT NULL
);`

	_, err = conn.Exec(ctx, sqlString)
	return err
}

func (m *migrator) Run(ctx context.Context, model database.ModelInterface) error {
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	migrationCreationTimes := make([]time.Time, len(migrations))
	iterator := 0
	for m := range migrations {
		migrationCreationTime, err := time.Parse(time.RFC3339, m)
		if err != nil {
			return errors.Wrap(err, "Failed to parse migration time.")
		}

		migrationCreationTimes[iterator] = migrationCreationTime
		iterator += 1
	}

	sort.Slice(migrationCreationTimes, func(i, j int) bool {
		return migrationCreationTimes[i].Before(migrationCreationTimes[j])
	})

	for i := 0; i < len(migrationCreationTimes); i = i + 1 {
		key := migrationCreationTimes[i].Format(time.RFC3339)
		tableMigrationKey := fmt.Sprintf(
			"%s_%s",
			key,
			model.Name(),
		)

		migrationHasRun, err := m.HasMigration(ctx, conn, tableMigrationKey)
		if err != nil {
			return errors.Wrap(err, "Failed to fetch migration.")
		}

		if migrationHasRun {
			continue
		}

		migration := migrations[key]
		if err = migration.Up(ctx, conn, model); err != nil {
			return errors.Wrap(err, "Failed to run migration.")
		}

		if err = m.RecordMigration(ctx, conn, tableMigrationKey); err != nil {
			return errors.Wrap(err, "Failed to record migration.")
		}
	}

	return nil
}

func (m *migrator) HasMigration(ctx context.Context, conn *pgxpool.Conn, identifier string) (bool, error) {
	count := 0
	err := conn.QueryRow(
		ctx,
		fmt.Sprintf(`SELECT COUNT(*) FROM migrations where id = $1`),
		identifier,
	).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *migrator) RecordMigration(ctx context.Context, conn *pgxpool.Conn, identifier string) error {
	query := `INSERT INTO migrations (id, executed_at) VALUES($1, CURRENT_TIMESTAMP)`
	_, err := conn.Exec(
		ctx,
		query,
		identifier,
	)

	if err != nil {
		return errors.Wrap(err, "Failed to record migration execution.")
	}

	return nil
}

func (m *migrator) DeleteMigration(ctx context.Context, conn *pgxpool.Conn, identifier string) error {
	_, err := conn.Exec(
		ctx,
		`DELETE FROM migrations WHERE id = $1`,
		identifier,
	)

	if err != nil {
		return errors.Wrap(err, "Failed to delete migration.")
	}

	return nil
}
