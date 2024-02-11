package content

import (
	"fmt"
	"github.com/fanky5g/ponzu/internal/domain/entities"
)

func (s *service) GetContent(entityType, entityId string) (interface{}, error) {
	target := fmt.Sprintf("%s:%s", entityType, entityId)
	return s.repository.FindOneByTarget(target)
}

func (s *service) GetContentBySlug(slug string) (string, interface{}, error) {
	return s.repository.FindOneBySlug(slug)
}

func (s *service) GetAll(namespace string) ([]interface{}, error) {
	return s.repository.FindAll(namespace)
}

func (s *service) GetAllWithOptions(namespace string, search *entities.Search) (int, []interface{}, error) {
	return s.repository.FindAllWithOptions(namespace, search.SortOrder, search.Pagination)
}
