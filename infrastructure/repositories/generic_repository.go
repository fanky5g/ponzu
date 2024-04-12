package repositories

import "github.com/fanky5g/ponzu/entities"

type GenericRepositoryInterface interface {
	Insert(entity interface{}) (interface{}, error)
	Latest() (interface{}, error)
	UpdateById(id string, update interface{}) (interface{}, error)
	Search(search *entities.Search) (int, []interface{}, error)
	FindOneById(id string) (interface{}, error)
	FindOneBy(field string, value interface{}) (interface{}, error)
	FindAll() ([]interface{}, error)
	DeleteById(id string) error
	//BatchInsert(entities []interface{}) ([]interface{}, error)
	//FindBy(criteria string, value interface{}) ([]interface{}, error)
	//FindAllOrderBy(
	//	order constants.SortOrder,
	//) ([]interface{}, error)
	//DeleteByTimeRange(start, end time.Time) error
}
