package storage

import (
	"fmt"

	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/pkg/errors"
)

func (s *service) DeleteFile(entityIds ...string) error {
	for _, entityId := range entityIds {
		f, err := s.GetFileUpload(entityId)
		if err != nil {
			return err
		}

		if f == nil {
			return nil
		}

		if err = s.client.Delete(f.Path); err != nil {
			return fmt.Errorf("failed to delete from file store: %v", err)
		}

		if err = s.repository.DeleteById(entityId); err != nil {
			return errors.Wrap(err, "Failed to delete item from database")
		}

		if err = s.searchClient.Delete(constants.UploadsEntityName, entityId); err != nil {
			return errors.Wrap(err, "Failed to delete search index entry")
		}
	}

	return nil
}
