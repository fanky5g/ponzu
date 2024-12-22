package datasource

import "io"

type Datasource interface {
	GetContentDisposition() string
	GetContentType() string
	io.Reader
}

type ChunkedData interface {
	GetCount() (int, error)
	GetChunkSize() int
	LoadChunk(size int, offset int) ([]interface{}, error)
}

type RowFormatter interface {
	FormatRow(interface{}) (interface{}, error)
}
