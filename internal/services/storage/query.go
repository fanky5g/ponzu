package storage

import (
	"github.com/fanky5g/ponzu/entities"
)

func (s *service) GetFileUpload(entityId string) (*entities.FileUpload, error) {
	file, err := s.repository.FindOneById(entityId)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, nil
	}

	return file.(*entities.FileUpload), nil
}

func (s *service) GetAllWithOptions(search *entities.Search) (int, []*entities.FileUpload, error) {
	total, files, err := s.repository.Find(search.SortOrder, search.Pagination)
	if err != nil {
		return 0, nil, err
	}

	if len(files) > 0 {
		out := make([]*entities.FileUpload, len(files))
		for i := range files {
			out[i] = files[i].(*entities.FileUpload)
		}

		return total, out, nil
	}

	return total, nil, nil
}
