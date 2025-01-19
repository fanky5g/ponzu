package mocks

import (
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	*mock.Mock
}

func (repo *Repository) Insert(entity interface{}) (interface{}, error) {
	args := repo.Called(entity)
	return args.Get(0), args.Error(1)
}

func (repo *Repository) Latest() (interface{}, error) {
	args := repo.Called()
	return args.Get(0), args.Error(1)
}

func (repo *Repository) UpdateById(id string, update interface{}) (interface{}, error) {
	args := repo.Called(id, update)
	return args.Get(0), args.Error(1)
}

func (repo *Repository) GetNumberOfRows() (int, error) {
	args := repo.Called()
	return args.Int(0), args.Error(1)
}

func (repo *Repository) Find(order constants.SortOrder, count, offset int) (int, []interface{}, error) {
	args := repo.Called(order, count, offset)

	result := args.Get(1)
	if result != nil {
		return args.Int(0), result.([]interface{}), args.Error(2)
	}

	return args.Int(0), nil, args.Error(2)
}

func (repo *Repository) FindOneById(id string) (interface{}, error) {
	args := repo.Called(id)
	return args.Get(0), args.Error(1)
}

func (repo *Repository) FindOneBy(criteria map[string]interface{}) (interface{}, error) {
	args := repo.Called(criteria)
	return args.Get(0), args.Error(1)
}

func (repo *Repository) FindAll() ([]interface{}, error) {
	args := repo.Called()

	result := args.Get(0)
	if result != nil {
		return result.([]interface{}), args.Error(1)
	}

	return nil, args.Error(1)
}

func (repo *Repository) FindByIds(ids ...string) ([]interface{}, error) {
	args := repo.Called(ids)

	result := args.Get(1)
	if result != nil {
		return nil, result.(error)
	}

	return args.Get(0).([]interface{}), args.Error(1)
}

func (repo *Repository) DeleteById(id string) error {
	args := repo.Called(id)
	return args.Error(0)
}

func (repo *Repository) DeleteBy(field string, operator constants.ComparisonOperator, value interface{}) error {
	args := repo.Called(field, operator, value)
	return args.Error(0)
}
