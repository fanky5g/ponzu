package formatter

import (
	"encoding/json"
	"errors"

	"github.com/fanky5g/ponzu/internal/datasource"
	"github.com/tidwall/gjson"
)

var ErrInvalidRow = errors.New("not a valid row")

type jsonFormatter struct {}

func (f *jsonFormatter) FormatRow(row interface{}) (interface{}, error) {
	rowType, ok := row.(datasource.Row)
	if !ok {
		return nil, ErrInvalidRow 
	}

	rowBuf := make([]string, 0)
	var rowData []byte

	rowData, err := json.Marshal(row)
	if err != nil {
		return nil, err
	}

	for _, col := range rowType.Columns() {
		result := gjson.GetBytes(rowData, col)
		rowBuf = append(rowBuf, result.String())
	}

	return rowBuf, nil
}

func NewJSONRowFormatter() (*jsonFormatter, error) {
	return &jsonFormatter{}, nil
}
