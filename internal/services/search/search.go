package search

import (
	"fmt"
	"github.com/fanky5g/ponzu/driver"
)

type service struct {
	client driver.SearchClientInterface
}

type Service interface {
	Search(entityName, query string, count, offset int) ([]interface{}, error)
}

func New(client driver.SearchClientInterface) (Service, error) {
	return &service{client: client}, nil
}

// Search conducts a search and returns a set of Ponzu "targets", Type:ID pairs,
// and an error. If there is no search index for the typeName (Type) provided,
// db.ErrNoIndex will be returned as the error
func (s *service) Search(entityName, query string, count, offset int) ([]interface{}, error) {
	index, err := s.client.GetIndex(entityName)
	if err != nil {
		return nil, fmt.Errorf("failed to get index for entity: %s", entityName)
	}

	return index.Search(query, count, offset)
}
