package content

import (
	"fmt"

	"github.com/fanky5g/ponzu/constants"
    "github.com/pkg/errors"	
)

func (s *service) DeleteContent(entityType, entityId string) error {
	repository := s.repository(entityType)
	if err := repository.DeleteById(entityId); err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}

	if err := s.slugRepository.DeleteBy("entity_id", constants.Equal, entityId); err != nil {
		return fmt.Errorf("failed to delete slug: %v", err)
	}

    if err := s.searchClient.Delete(entityType, entityId); err != nil {
        return errors.Wrap(err, "failed to delete indexed content")
    }

	return nil
}
