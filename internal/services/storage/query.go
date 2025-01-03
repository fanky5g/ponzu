package storage

import (
	contentEntities "github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/internal/search"
)

func (s *service) GetUpload(entityId string) (*contentEntities.Upload, error) {
	file, err := s.repository.FindOneById(entityId)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, nil
	}

	return file.(*contentEntities.Upload), nil
}

func (s *service) GetAllWithOptions(search *search.Search) (int, []*contentEntities.Upload, error) {
	total, files, err := s.repository.Find(search.SortOrder, search.Count, search.Offset)
	if err != nil {
		return 0, nil, err
	}

	if len(files) > 0 {
		out := make([]*contentEntities.Upload, len(files))
		for i := range files {
			out[i] = files[i].(*contentEntities.Upload)
		}

		return total, out, nil
	}

	return total, nil, nil
}
