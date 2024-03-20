package storage

import (
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/entities"
)

func (s *service) GetFileUpload(key string) (*entities.FileUpload, error) {
	file, err := s.Service.GetContent(constants.UploadsEntityName, key)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, nil
	}

	return file.(*entities.FileUpload), nil
}

func (s *service) GetAllUploads() ([]entities.FileUpload, error) {
	files, err := s.Service.GetAll(constants.UploadsEntityName)
	if err != nil {
		return nil, err
	}

	f := make([]entities.FileUpload, len(files))
	for i, file := range files {
		f[i] = file.(entities.FileUpload)
	}

	return f, nil
}
