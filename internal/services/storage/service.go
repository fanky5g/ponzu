// Package storage provides a re-usable file storage and storage utility for Ponzu
// systems to handle multipart form data.
package storage

import (
	contentEntities "github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/search"
	"github.com/fanky5g/ponzu/tokens"
	"mime/multipart"
)

type service struct {
	client       driver.StorageClientInterface
	repository   driver.Repository
	searchClient driver.SearchInterface
}

type Service interface {
	GetAllWithOptions(search *search.Search) (int, []*contentEntities.FileUpload, error)
	GetFileUpload(target string) (*contentEntities.FileUpload, error)
	DeleteFile(target ...string) error
	StoreFiles(files map[string]*multipart.FileHeader) (map[string]string, error)
	driver.StaticFileSystemInterface
}

func New(
	db driver.Database,
	searchClient driver.SearchInterface,
	client driver.StorageClientInterface) (Service, error) {
	s := &service{
		client:       client,
		searchClient: searchClient,
		repository:   db.GetRepositoryByToken(tokens.UploadRepositoryToken),
	}

	return s, nil
}
