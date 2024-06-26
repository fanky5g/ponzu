package content

import (
	"fmt"
	"github.com/fanky5g/ponzu/content/item"
)

// UpdateContent supports only full updates. The entire structure will be overwritten.
// Partial updates are not supported.
func (s *service) UpdateContent(entityType, entityId string, update interface{}) (interface{}, error) {
	identifiable, ok := update.(item.Identifiable)
	if !ok {
		return nil, fmt.Errorf("update not supported for %s", entityType)
	}

	identifiable.SetItemID(entityId)
	var sluggable item.Sluggable
	if sluggable, ok = update.(item.Sluggable); ok {
		slug, err := s.GetSlug(entityType, entityId)
		if err != nil {
			return nil, err
		}

		if slug != nil {
			sluggable.SetSlug(slug.Slug)
		}
	}

	return s.repository(entityType).UpdateById(entityId, update)
}
