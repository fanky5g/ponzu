package content

import (
	"github.com/fanky5g/ponzu/entities"
)

func (s *service) GetContent(entityType, entityId string) (interface{}, error) {
	return s.repository(entityType).FindOneById(entityId)
}

func (s *service) GetContentBySlug(slug string) (interface{}, error) {
	match, err := s.slugRepository.FindOneBy(map[string]interface{}{
		"slug": slug,
	})

	if err != nil {
		return nil, err
	}

	if match == nil {
		return nil, nil
	}

	ss := match.(*entities.Slug)
	return s.repository(ss.EntityType).FindOneById(ss.EntityId)
}

func (s *service) GetAll(entityType string) ([]interface{}, error) {
	return s.repository(entityType).FindAll()
}

func (s *service) GetAllWithOptions(entityType string, search *entities.Search) (int, []interface{}, error) {
	return s.repository(entityType).Search(search)
}
