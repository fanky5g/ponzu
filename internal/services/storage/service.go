// Package storage provides a re-usable file storage and storage utility for Ponzu
// systems to handle multipart form data.
package storage

import (
	"mime/multipart"

	contentEntities "github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/internal/search"
	"github.com/fanky5g/ponzu/tokens"
)

type service struct {
	client       driver.StorageClientInterface
	repository   database.Repository
	searchClient search.SearchInterface
}

type Service interface {
	GetAllWithOptions(search *search.Search) (int, []*contentEntities.Upload, error)
	GetUpload(target string) (*contentEntities.Upload, error)
	DeleteUpload(target ...string) error
	UploadFiles(files map[string]*multipart.FileHeader) (map[string]string, error)
	driver.StaticFileSystemInterface
}

func New(
	db database.Database,
	searchClient search.SearchInterface,
	client driver.StorageClientInterface) (Service, error) {
	s := &service{
		client:       client,
		searchClient: searchClient,
		repository:   db.GetRepositoryByToken(tokens.UploadRepositoryToken),
	}

	return s, nil
}
