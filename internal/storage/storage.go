package storage

import (
	"io"
	"net/http"
)

type Client interface {
	Save(fileName string, file io.ReadCloser) (string, int64, error)
	Delete(path string) error
	Open(name string) (http.File, error)
	Attributes(name string) (*FileAttributes, error)
}

type FileAttributes struct {
	ContentType string
}
