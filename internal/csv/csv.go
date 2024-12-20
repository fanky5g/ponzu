package csv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fanky5g/ponzu/internal/datasource"
	"github.com/tidwall/gjson"
)

type csvFile struct {
	b          *bytes.Buffer
	w          *csv.Writer
	header     []string
	read       int
	totalRows  int
	datasource datasource.DataSource
}

func New(dataSource datasource.DataSource) (io.Reader, error) {
	b := new(bytes.Buffer)
	w := csv.NewWriter(bufio.NewWriter(b))

	f := &csvFile{
		b: b,
		w: w,
	}

	var err error
	if f.header, err = dataSource.GetColumns(); err != nil {
		return nil, fmt.Errorf("failed to get column headers: %v", err)
	}

	if f.totalRows, err = dataSource.GetNumberOfRows(); err != nil {
		return nil, fmt.Errorf("failed to get number of rows: %v", err)
	}

	if err := w.Write(f.header); err != nil {
		return nil, fmt.Errorf("failed to write column headers: %v", err)
	}

	return f, nil
}

func (f *csvFile) Read(p []byte) (int, error) {
	n, err := f.b.Read(p)
	if err != nil {
		if err == io.EOF && f.read != f.totalRows {
			var data []interface{}
			data, err = f.datasource.LoadData(f.read)
			if err != nil {
				return 0, err
			}

			if err = f.write(data); err != nil {
				return 0, err
			}

			f.read += len(data)
			return f.b.Read(p)
		}
	}

	return n, err
}

func (f *csvFile) write(data []interface{}) error {
	var err error
	for row := range data {
		rowBuf := make([]string, 0)
		var rowData []byte
		rowData, err = json.Marshal(data[row])
		if err != nil {
			break
		}

		for _, col := range f.header {
			result := gjson.GetBytes(rowData, col)
			rowBuf = append(rowBuf, result.String())
		}

		if err = f.w.Write(rowBuf); err != nil {
			break
		}
	}

	if err != nil {
		return err
	}

	f.w.Flush()
	return nil
}
