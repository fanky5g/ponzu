package content

import (
	"github.com/fanky5g/ponzu/internal/search"
)

var UploadType = "Upload"

func (s *Service) ListReferences(typeName string, searchQuery *search.Search) ([]interface{}, int, error) {
	switch typeName {
	case UploadType:
		return s.uploadService.GetAllWithOptions(searchQuery)
	default:
		return s.GetAllWithOptions(typeName, searchQuery)
	}
}

func (s *Service) GetReference(entityName, entityId string) (interface{}, error) {
	switch entityName {
	case UploadType:
		return s.uploadService.GetUpload(entityId)
	default:
		return s.GetContent(entityName, entityId)
	}
}

func (s *Service) GetReferences(entityName string, entityIds ...string) ([]interface{}, error) {
	switch entityName {
	case UploadType:
		ups, err := s.uploadService.GetUploads(entityIds...)
		if err != nil {
			return nil, err
		}

		out := make([]interface{}, len(ups))
		for i, up := range ups {
			out[i] = up
		}

		return out, nil
	default:
		return s.GetContentByIds(entityName, entityIds...)
	}
}
