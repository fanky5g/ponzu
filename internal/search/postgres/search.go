package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/internal/database/postgres/connection"
	"github.com/fanky5g/ponzu/internal/database/postgres/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidSearchEntity = errors.New("invalid search entity")
)

func (c *Client) Search(entity interface{}, query string, limit, offset int) ([]interface{}, error) {
	results, _, err := c.SearchWithPagination(entity, query, limit, offset)
	return results, err
}

func (c *Client) SearchWithPagination(entity interface{}, query string, limit, offset int) ([]interface{}, int, error) {
	persistable, ok := entity.(database.Persistable)
	if !ok {
		return nil, 0, ErrInvalidSearchEntity
	}

	searchableFields, err := getSearchableFields(entity)
	if err != nil {
		return nil, 0, err
	}

	ctx := context.Background()
	pool, err := connection.Get(ctx)
	if err != nil {
		return nil, 0, err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, 0, err
	}

	defer conn.Release()
	databaseRepository := c.db.GetRepositoryByToken(persistable.GetRepositoryToken())
	if databaseRepository == nil {
		return nil, 0, ErrInvalidSearchEntity
	}

	repo, ok := databaseRepository.(*repository.Repository)
	if !ok {
		return nil, 0, ErrInvalidSearchEntity
	}

	queryLength := len(searchableFields)
	whereClauses := make([]string, queryLength)
	values := make([]interface{}, queryLength)
	position := 0
	for _, field := range searchableFields {
		whereClauses[position] = fmt.Sprintf(
			"(document->>'%s') ILIKE $1",
			field,
		)

		values[position] = query
		position = position + 1
	}

	size := repository.DefaultQuerySize
	if limit > 0 {
		size = limit
	}

	whereClause := strings.Join(whereClauses, " OR ")

	sqlString := fmt.Sprintf(`
			SELECT id, created_at, updated_at, document
			FROM %s
			WHERE (%s) AND deleted_at IS NULL
			ORDER BY updated_at DESC
			LIMIT %d
	`, repo.TableName(), whereClause, size)

	if offset > 0 {
		sqlString = fmt.Sprintf(`
			%s
			OFFSET %d
`, sqlString, offset)
	}

	value := fmt.Sprintf("%%%s%%", query)
	count, err := c.count(ctx, conn, repo.TableName(), whereClause, value)
	if err != nil {
		return nil, 0, err
	}

	rows, err := conn.Query(ctx, sqlString, value)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	results := make([]interface{}, 0)
	for rows.Next() {
		var result interface{}
		if result, err = repo.ScanRow(rows); err != nil {
			return nil, 0, err
		}

		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func (c *Client) count(ctx context.Context, conn *pgxpool.Conn, tableName, whereClause, value string) (int, error) {
	sqlString := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE (%s) AND deleted_at IS NULL
`, tableName, whereClause)

	count := 0
	err := conn.QueryRow(ctx, sqlString, value).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Update is a no-op as with postgres we don't have to actually update any index
func (c *Client) Update(id string, data interface{}) error {
	return nil
}

// Delete is a no-op as with postgres we don't have to delete from an index.
func (c *Client) Delete(entityName, id string) error {
	return nil
}
