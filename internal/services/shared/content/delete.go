package content

import (
	"fmt"
	"log"
)

func (s *service) DeleteContent(entityType, entityId string) error {
	target := fmt.Sprintf("%s:%s", entityType, entityId)
	if err := s.repository.DeleteEntity(target); err != nil {
		return err
	}

	index, err := s.searchClient.GetIndex(s.getEntityType(target))
	if err != nil {
		log.Printf("failed to delete search index: %v", err)
	}

	if index != nil {
		return index.Delete(target)
	}

	return nil
}
