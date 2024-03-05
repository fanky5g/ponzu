package storage

import (
	"github.com/fanky5g/ponzu/internal/domain/entities"
)

func (s *service) GetFileUpload(key string) (*entities.FileUpload, error) {
	file, err := s.Service.GetContent(UploadsEntityName, key)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, nil
	}

	return file.(*entities.FileUpload), nil
}

func (s *service) GetAllUploads() ([]entities.FileUpload, error) {
	files, err := s.Service.GetAll(UploadsEntityName)
	if err != nil {
		return nil, err
	}

	f := make([]entities.FileUpload, len(files))
	for i, file := range files {
		f[i] = file.(entities.FileUpload)
	}

	return f, nil
}
