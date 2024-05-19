package driver

type SearchInterface interface {
	Update(id string, data interface{}) error
	Delete(id string) error
	Search(query string, count, offset int) ([]interface{}, int, error)
}

type SearchInterface interface {
	SearchIndexInterface
	SearchWithPagination(query string, count, offset int) ([]interface{}, int, error)
}
