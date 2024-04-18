// Package storage provides a re-usable file storage and storage utility for Ponzu
// systems to handle multipart form data.
package storage

import (
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/models"
	"github.com/fanky5g/ponzu/tokens"
	"mime/multipart"
)

type service struct {
	client     driver.StorageClientInterface
	repository driver.Repository
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
	searchClient driver.SearchClientInterface,
	client driver.StorageClientInterface) (Service, error) {

	if err := searchClient.CreateIndex(
		models.WrapPonzuModelNameSpace(tokens.Repository(constants.UploadsEntityName)), &entities.FileUpload{}); err != nil {
		return nil, err
	}

	s := &service{
		client: client,
		repository: db.Get(
			models.WrapPonzuModelNameSpace(tokens.UploadRepositoryToken),
		).(driver.Repository),
	}

	return s, nil
}
