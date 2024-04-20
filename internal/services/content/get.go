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

func (s *service) GetSlug(entityType, entityId string) (*entities.Slug, error) {
	slug, err := s.slugRepository.FindOneBy(map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityId,
	})

	if err != nil {
		return nil, err
	}

	if slug == nil {
		return nil, nil
	}

	return slug.(*entities.Slug), nil
}

func (s *service) GetAll(entityType string) ([]interface{}, error) {
	return s.repository(entityType).FindAll()
}

func (s *service) GetAllWithOptions(entityType string, search *entities.Search) (int, []interface{}, error) {
	return s.repository(entityType).Find(search.SortOrder, search.Pagination)
}
