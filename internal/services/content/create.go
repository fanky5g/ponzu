package content

import (
	"strings"
)

func (s *service) CreateContent(entityType string, entity interface{}) (string, error) {
	//identifiable, ok := entity.(item.Identifiable)
	//if !ok {
	//	return "", errors.New("item does not implement identifiable interface")
	//}
	//
	//// TODO: setting unique ids should be up to the database layer
	//if identifiable.UniqueID().IsNil() {
	//	// add UUID to data for use in embedded Item
	//	uid, err := uuid.NewV4()
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	entity.(item.Identifiable).SetUniqueID(uid)
	//}
	//
	//// TODO: setting slugs should be up to the database layer - or a slug service
	//if sluggable, ok := entity.(item.Sluggable); ok && sluggable.ItemSlug() == "" {
	//	slug, err := item.Slug(entity.(item.Identifiable))
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	slug, err = s.repository.UniqueSlug(slug)
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	entity.(item.Sluggable).SetSlug(slug)
	//}
	//
	//id, err := s.repository.SetEntity(entityType, entity)
	//if err != nil {
	//	return "", err
	//}
	//
	//if searchable, ok := entity.(driver.Searchable); ok && searchable.IndexContent() {
	//	var index driver.SearchIndexInterface
	//	index, err = s.searchClient.GetIndex(s.getEntityType(entityType))
	//	if err != nil {
	//		return "", fmt.Errorf("failed to index %s for search: %v", entityType, err)
	//	}
	//
	//	if err = index.Update(id, entity); err != nil {
	//		return "", err
	//	}
	//}
	//
	//return fmt.Sprint(id), nil
	return "", nil
}

func (s *service) getEntityType(target string) string {
	return strings.Split(target, ":")[0]
}
