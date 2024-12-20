package database

import (
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/search"
)

type Repository interface {
	Insert(entity interface{}) (interface{}, error)
	Latest() (interface{}, error)
	UpdateById(id string, update interface{}) (interface{}, error)
	Find(order constants.SortOrder, pagination *search.Pagination) (int, []interface{}, error)
	FindOneById(id string) (interface{}, error)
	FindOneBy(criteria map[string]interface{}) (interface{}, error)
	FindAll() ([]interface{}, error)
	DeleteById(id string) error
	DeleteBy(field string, operator constants.ComparisonOperator, value interface{}) error
}
