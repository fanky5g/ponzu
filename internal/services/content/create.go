package content

import (
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/util"
)

func (s *service) CreateContent(entityType string, entity interface{}) (string, error) {
	repository := s.repository(entityType)
	identifiable, ok := entity.(item.Identifiable)
	if !ok {
		return "", errors.New("item does not implement identifiable interface")
	}

	// if slug is empty from the request - slugify the name/title of the entity. Slugs are unique identifiers.
	// Non-unique slugs will throw an error. It's up to the content manager or API to send a unique slug identifier.
	// Currently, we are not sending useful messages to the client, but in the future, we must send a useful message
	// that informs the content manager to update the slug. This will be done after *system errors are implemented.
	// after creation slugs cannot be updated.
	sluggable, ok := entity.(item.Sluggable)
	if !ok {
		return "", errors.New("entity does not implement sluggable interface")
	}

	if sluggable.ItemSlug() == "" {
		slug, err := util.Slugify(sluggable.GetTitle())
		if err != nil {
			return "", fmt.Errorf("failed to get slug: %v", err)
		}

		sluggable.SetSlug(slug)
	}

	content, err := repository.Insert(entity)
	if err != nil {
		return "", fmt.Errorf("failed to create content: %v", err)
	}

	identifiable = content.(item.Identifiable)
	if _, err = s.slugRepository.Insert(&entities.Slug{
		EntityType: entityType,
		EntityId:   identifiable.ItemID(),
		Slug:       sluggable.ItemSlug(),
	}); err != nil {
		return "", fmt.Errorf("failed to save slug: %v", err)
	}

	if searchable, ok := entity.(entities.Searchable); ok && searchable.IndexContent() {
		var index driver.SearchIndexInterface
		index, err = s.searchClient.GetIndex(entityType)
		if err != nil {
			return "", fmt.Errorf("failed to get search index for %v: %v", entityType, err)
		}

		if err = index.Update(identifiable.ItemID(), entity); err != nil {
			return "", fmt.Errorf("failed to index entity: %v", err)
		}
	}

	return identifiable.ItemID(), nil
}
