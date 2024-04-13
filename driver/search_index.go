package driver

type SearchIndexInterface interface {
	Update(id string, data interface{}) error
	Delete(id string) error
	Search(query string, count, offset int) ([]interface{}, error)
}
