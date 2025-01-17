package database

import (
	"github.com/fanky5g/ponzu/internal/constants"
)

type Repository interface {
	Insert(entity interface{}) (interface{}, error)
	Latest() (interface{}, error)
	UpdateById(id string, update interface{}) (interface{}, error)
	GetNumberOfRows() (int, error)
	Find(order constants.SortOrder, count, offset int) (int, []interface{}, error)
	FindOneById(id string) (interface{}, error)
	FindOneBy(criteria map[string]interface{}) (interface{}, error)
	FindAll() ([]interface{}, error)
	FindByIds(ids ...string) ([]interface{}, error)
	DeleteById(id string) error
	DeleteBy(field string, operator constants.ComparisonOperator, value interface{}) error
}
