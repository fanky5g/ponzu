package dataexporter

import (
	"errors"

	"github.com/fanky5g/ponzu/internal/content/dataexporter/csv"
	"github.com/fanky5g/ponzu/internal/datasource"
)

var ErrUnsupportedDataType = errors.New("unsupported export datatype")

type DataExporter interface {
	Export(kind, fileName string, data datasource.ChunkedData) (datasource.Datasource, error)
}

type exporter struct {
	rowFormatter datasource.RowFormatter
}

func (e *exporter) Export(kind, fileName string, data datasource.ChunkedData) (datasource.Datasource, error) {
	switch kind {
	case "csv":
		return csv.New(fileName, data, e.rowFormatter)
	default:
		return nil, ErrUnsupportedDataType
	}
}

func New(rowFormatter datasource.RowFormatter) (*exporter, error) {
	return &exporter{rowFormatter: rowFormatter}, nil
}
