package repository

import (
	"context"
	"fmt"
	"github.com/fanky5g/ponzu/internal/constants"
)

func (repo *Repository) DeleteById(id string) error {
	sqlString := fmt.Sprintf(
		"UPDATE %s SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1::uuid",
		repo.model.Name(),
	)

	ctx := context.Background()
	conn, err := repo.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()
	_, err = conn.Exec(ctx, sqlString, id)

	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}

	return nil
}

func (repo *Repository) DeleteBy(field string, operator constants.ComparisonOperator, value interface{}) error {
	valueType, err := repo.valueType(value)
	if err != nil {
		return err
	}

	comparisonOperator, err := repo.getComparisonOperator(operator)
	if err != nil {
		return err
	}

	sqlString := fmt.Sprintf(
		"UPDATE %s SET deleted_at = CURRENT_TIMESTAMP WHERE (document->>'%s')::%s %s $1::%s",
		repo.model.Name(),
		field,
		valueType,
		comparisonOperator,
		valueType,
	)

	ctx := context.Background()
	conn, err := repo.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()
	_, err = conn.Exec(ctx, sqlString, value)

	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}

	return nil
}
