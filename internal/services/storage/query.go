package storage

import (
	contentEntities "github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/search"
)

func (s *service) GetFileUpload(entityId string) (*contentEntities.FileUpload, error) {
	file, err := s.repository.FindOneById(entityId)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, nil
	}

	return file.(*contentEntities.FileUpload), nil
}

func (s *service) GetAllWithOptions(search *search.Search) (int, []*contentEntities.FileUpload, error) {
	total, files, err := s.repository.Find(search.SortOrder, search.Pagination)
	if err != nil {
		return 0, nil, err
	}

	if len(files) > 0 {
		out := make([]*contentEntities.FileUpload, len(files))
		for i := range files {
			out[i] = files[i].(*contentEntities.FileUpload)
		}

		return total, out, nil
	}

	return total, nil, nil
}
