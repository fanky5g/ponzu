package content

import (
	"fmt"
	"github.com/fanky5g/ponzu/constants"
	log "github.com/sirupsen/logrus"
)

func (s *service) DeleteContent(entityType, entityId string) error {
	repository := s.repository(entityType)
	if err := repository.DeleteById(entityId); err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}

	if err := s.slugRepository.DeleteBy("entity_id", constants.Equal, entityId); err != nil {
		return fmt.Errorf("failed to delete slug: %v", err)
	}

	index, err := s.searchClient.GetIndex(entityType)
	if err != nil {
		log.WithFields(log.Fields{
			"Error":      err,
			"EntityType": entityType,
		}).Warning("Failed to delete search index")
	}

	if index != nil {
		return index.Delete(entityId)
	}

	return nil
}
