// Package storage provides a re-usable file storage and storage utility for Ponzu
// systems to handle multipart form data.
package storage

import (
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/tokens"
	"mime/multipart"
)

type service struct {
	client       driver.StorageClientInterface
	repository   driver.Repository
	searchClient driver.SearchInterface
}

type Service interface {
	GetAllWithOptions(search *entities.Search) (int, []*entities.FileUpload, error)
	GetFileUpload(target string) (*entities.FileUpload, error)
	DeleteFile(target string) error
	StoreFiles(files map[string]*multipart.FileHeader) (map[string]string, error)
	driver.StaticFileSystemInterface
}

func New(
	db driver.Database,
	searchClient driver.SearchInterface,
	client driver.StorageClientInterface) (Service, error) {
	s := &service{
		client:     client,
		repository: db.GetRepositoryByToken(tokens.UploadRepositoryToken),
	}

	return s, nil
}
