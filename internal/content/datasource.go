package content

import (
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/entities"
)

var chunkSize = 50

type contentDatasource struct {
	entity         interface{}
	entityName     string
	contentService *Service
}

func (d *contentDatasource) GetCount() (int, error) {
	return d.contentService.GetNumberOfRows(d.entityName)
}

func (d *contentDatasource) GetChunkSize() int {
	return chunkSize
}

func (d *contentDatasource) LoadChunk(size, offset int) ([]interface{}, error) {
	chunk, _, e := d.contentService.GetAllWithOptions(d.entityName, &entities.Search{
		SortOrder: constants.Descending,
		Pagination: &entities.Pagination{
			Count:  size,
			Offset: offset,
		},
	})

	return chunk, e
}
