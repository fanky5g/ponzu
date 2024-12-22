package csv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/fanky5g/ponzu/internal/datasource"
)

var (
	ErrInvalidRowData = errors.New("invalid row data. must be []string")
	ErrInvalidRow     = errors.New("invalid row")
)

type csvData struct {
	b            *bytes.Buffer
	w            *csv.Writer
	read         int
	datasource   datasource.ChunkedData
	rowFormatter datasource.RowFormatter
	fileName     string
	rowCount     int
}

func New(fileName string, datasource datasource.ChunkedData, rowFormatter datasource.RowFormatter) (*csvData, error) {
	b := new(bytes.Buffer)
	w := csv.NewWriter(bufio.NewWriter(b))

	rowCount, err := datasource.GetCount()
	if err != nil {
		return nil, err
	}

	f := &csvData{
		b:            b,
		w:            w,
		datasource:   datasource,
		fileName:     fileName,
		rowCount:     rowCount,
		rowFormatter: rowFormatter,
	}

	return f, nil
}

func (f *csvData) GetContentDisposition() string {
	return fmt.Sprintf(`attachment; filename="export-%s-%d.csv"`, f.fileName, time.Now().Unix())
}

func (f *csvData) GetContentType() string {
	return "text/csv"
}

func (f *csvData) Read(p []byte) (int, error) {
	if f.rowCount == 0 {
		return 0, nil
	}

	n, err := f.b.Read(p)
	if err != nil {
		if err == io.EOF && f.read != f.rowCount {
			var data []interface{}
			data, err = f.datasource.LoadChunk(f.datasource.GetChunkSize(), f.read)
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

func (f *csvData) write(data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	if f.read == 0 {
		row, ok := data[0].(datasource.Row)
		if !ok {
			return ErrInvalidRow
		}

		if err := f.w.Write(row.Columns()); err != nil {
			return err
		}
	}

	for i := range data {
		formattedRow, err := f.rowFormatter.FormatRow(data[i])
		if err != nil {
			return err
		}

		row, ok := formattedRow.([]string)
		if !ok {
			return ErrInvalidRow
		}

		if err = f.w.Write(row); err != nil {
			break
		}
	}

	f.w.Flush()
	return nil
}
