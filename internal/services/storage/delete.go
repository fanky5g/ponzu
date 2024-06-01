package storage

import (
	"fmt"
)

func (s *service) DeleteFile(entityId string) error {
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

	return s.repository.DeleteById(entityId)
}
