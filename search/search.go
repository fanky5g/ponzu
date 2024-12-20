package search

type Client interface {
	Update(id string, data interface{}) error
	Delete(entityName, entityId string) error
	Search(entityDefinition interface{}, query string, count, offset int) ([]interface{}, error)
	SearchWithPagination(entityDefinition interface{}, query string, count, offset int) ([]interface{}, int, error)
}
