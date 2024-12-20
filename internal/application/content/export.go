package content

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/entities"
	"github.com/fanky5g/ponzu/search"
	"github.com/tidwall/gjson"
)

var chunkSize = 50

type csvFile struct {
	b         *bytes.Buffer
	totalRows int
	fields    []string
	read      int
	loadMore  func(offset int) ([]interface{}, error)
	w         *csv.Writer
}

func (f *csvFile) Read(p []byte) (int, error) {
	n, err := f.b.Read(p)
	if err != nil {
		if err == io.EOF && f.read != f.totalRows {
			var data []interface{}
			data, err = f.loadMore(f.read)
			if err != nil {
				return 0, err
			}

			if err = f.WriteJSONData(data); err != nil {
				return 0, err
			}

			f.read += len(data)
			return f.b.Read(p)
		}
	}

	return n, err
}

func (f *csvFile) WriteJSONData(data []interface{}) error {
	var err error
	for row := range data {
		rowBuf := make([]string, 0)
		var rowData []byte
		rowData, err = json.Marshal(data[row])
		if err != nil {
			break
		}

		for _, col := range f.fields {
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

func (s *Service) ExportCSV(entity interface{}) (*entities.ResponseStream, error) {
	entityInterface, ok := entity.(content.Entity)
	if !ok {
		return nil, errors.New("Invalid content type")
	}

	csvFormattable, ok := entity.(item.CSVFormattable)
	if !ok {
		return nil, fmt.Errorf("%s does not implement CSV Formattable interface", entityInterface.EntityName())
	}

	offset := 0
	var data []interface{}
	var size int
	var err error
	// get entities data and loop through creating a csv row per result
	data, size, err = s.GetAllWithOptions(entityInterface.EntityName(), &entities.Search{
		SortOrder: constants.Descending,
		Pagination: &search.Pagination{
			Count:  chunkSize,
			Offset: offset,
		},
	})

	if err != nil {
		return nil, err
	}

	if size == 0 {
		return nil, nil
	}

	fields := csvFormattable.FormatCSV()
	b := new(bytes.Buffer)
	w := csv.NewWriter(bufio.NewWriter(b))

	f := &csvFile{
		totalRows: size,
		fields:    fields,
		b:         b,
		w:         w,
		read:      len(data),
		loadMore: func(offset int) ([]interface{}, error) {
			d, _, e := s.GetAllWithOptions(entityInterface.EntityName(), &entities.Search{
				SortOrder: constants.Descending,
				Pagination: &search.Pagination{
					Count:  chunkSize,
					Offset: offset,
				},
			})

			return d, e
		},
	}

	if err = w.Write(fields); err != nil {
		return nil, fmt.Errorf("failed to write column headers: %v", err)
	}

	if err = f.WriteJSONData(data); err != nil {
		return nil, err
	}

	return &entities.ResponseStream{
		ContentType:        "text/csv",
		ContentDisposition: fmt.Sprintf(`attachment; filename="export-%s-%d.csv"`, entityInterface.EntityName(), time.Now().Unix()),
		Payload:            f,
	}, nil
}
