package references

import (
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/search"
	"github.com/fanky5g/ponzu/internal/uploads"
)

var UploadType = "Upload"

type Service struct {
	contentService *content.Service
	uploadService  *uploads.Service
}

func (s *Service) ListReferences(typeName string, searchQuery *search.Search) ([]interface{}, int, error) {
	switch typeName {
	case UploadType:
		return s.uploadService.GetAllWithOptions(searchQuery)
	default:
		return s.contentService.GetAllWithOptions(typeName, searchQuery)
	}
}

func (s *Service) GetReference(typeName, id string) (interface{}, error) {
	switch typeName {
	case UploadType:
		return s.uploadService.GetUpload(id)
	default:
		return s.contentService.GetContent(typeName, id)
	}
}

func New(contentService *content.Service, uploadService *uploads.Service) (*Service, error) {
	return &Service{
		contentService: contentService,
		uploadService:  uploadService,
	}, nil
}
