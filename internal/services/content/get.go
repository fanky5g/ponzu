package content

import (
	"github.com/fanky5g/ponzu/entities"
)

func (s *service) GetContent(entityType, entityId string) (interface{}, error) {
	//target := fmt.Sprintf("%s:%s", entityType, entityId)
	//return s.repository.FindOneByTarget(target)
	return nil, nil
}

func (s *service) GetContentBySlug(slug string) (string, interface{}, error) {
	//return s.repository.FindOneBySlug(slug)
	return "", nil, nil
}

func (s *service) GetAll(namespace string) ([]interface{}, error) {
	//return s.repository.FindAll(namespace)
	return nil, nil
}

func (s *service) GetAllWithOptions(namespace string, search *entities.Search) (int, []interface{}, error) {
	//return s.repository.FindAllWithOptions(namespace, search.SortOrder, search.Pagination)
	return 0, nil, nil
}
