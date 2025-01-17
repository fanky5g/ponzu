package references

import (
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/search"
	"github.com/fanky5g/ponzu/internal/uploads"
)

var UploadType = "Upload"

type Service struct {
	contentService *content.Service
	uploadService  *uploads.UploadService
}

func (s *Service) ListReferences(typeName string, searchQuery *search.Search) ([]interface{}, int, error) {
	switch typeName {
	case UploadType:
		return s.uploadService.GetAllWithOptions(searchQuery)
	default:
		return s.contentService.GetAllWithOptions(typeName, searchQuery)
	}
}

func (s *Service) GetReference(entityName, entityId string) (interface{}, error) {
	switch entityName {
	case UploadType:
		return s.uploadService.GetUpload(entityId)
	default:
		return s.contentService.GetContent(entityName, entityId)
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
		return s.contentService.GetContentByIds(entityName, entityIds...)
	}
}

func New(contentService *content.Service, uploadService *uploads.UploadService) (*Service, error) {
	return &Service{
		contentService: contentService,
		uploadService:  uploadService,
	}, nil
}
