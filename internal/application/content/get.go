package content

import (
	"github.com/fanky5g/ponzu/internal/entities"
)

func (s *Service) GetContent(entityType, entityId string) (interface{}, error) {
	return s.repository(entityType).FindOneById(entityId)
}

func (s *Service) GetContentBySlug(slug string) (interface{}, error) {
	match, err := s.slugs.FindOneBy(map[string]interface{}{
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

func (s *Service) GetSlug(entityType, entityId string) (*entities.Slug, error) {
	slug, err := s.slugs.FindOneBy(map[string]interface{}{
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

func (s *Service) GetAll(entityType string) ([]interface{}, error) {
	return s.repository(entityType).FindAll()
}

func (s *Service) GetAllWithOptions(entityType string, search *entities.Search) ([]interface{}, int, error) {
	count, matches, err := s.repository(entityType).Find(search.SortOrder, search.Pagination)
	return matches, count, err
}
