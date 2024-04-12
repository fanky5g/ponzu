package repositories

type GenericRepositoryInterface interface {
	Insert(entity interface{}) (interface{}, error)
	Latest() (interface{}, error)
	UpdateById(id string, update interface{}) (interface{}, error)
	//BatchInsert(entities []interface{}) ([]interface{}, error)
	//FindOneById(id string) (interface{}, error)
	//FindOneBy(criteria string, value interface{}) (interface{}, error)
	//FindBy(criteria string, value interface{}) ([]interface{}, error)
	//FindAll() ([]interface{}, error)
	//FindAllOrderBy(
	//	order constants.SortOrder,
	//) ([]interface{}, error)
	//FindAllPaginatedOrderBy(
	//	pagination *entities.Pagination,
	//	sortOrder constants.SortOrder,
	//)
	//DeleteById(id string) error
	//DeleteByTimeRange(start, end time.Time) error
}
