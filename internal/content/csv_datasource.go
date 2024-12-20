package content

import (
	"fmt"

	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/entities"
)

func (datasource *CSVExportDataSource) GetNumberOfRows() (int, error) {
	return datasource.contentService.GetNumberOfRows(datasource.entityName)
}

func (datasource *CSVExportDataSource) GetColumns() ([]string, error) {
	csvFormattable, ok := datasource.entity.(item.CSVFormattable)
	if !ok {
		return nil, fmt.Errorf("%s does not implement CSV Formattable interface", datasource.entityName)
	}

	return csvFormattable.FormatCSV(), nil
}

func (datasource *CSVExportDataSource) LoadData(offset int) ([]interface{}, error) {
	d, _, e := datasource.contentService.GetAllWithOptions(datasource.entityName, &entities.Search{
		SortOrder: constants.Descending,
		Pagination: &entities.Pagination{
			Count:  CSVChunkSize,
			Offset: offset,
		},
	})

	return d, e
}

func NewCSVExportDataSource(entity interface{}) (*CSVExportDataSource, error) {
	entityInterface, ok := entity.(content.Entity)
	if !ok {
		return nil, ErrInvalidContentType
	}

	return &CSVExportDataSource{
		entity:     entity,
		entityName: entityInterface.EntityName(),
	}, nil
}
