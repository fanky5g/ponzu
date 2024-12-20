package storage 

import (
	"net/http"
)

type FileSystemInterface interface {
	Open(name string) (http.File, error)
}
