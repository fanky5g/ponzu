// Package storage provides a re-usable file storage and storage utility for Ponzu
// systems to handle multipart form data.
package storage

import (
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
	"github.com/fanky5g/ponzu/internal/services/shared/content"
	"github.com/fanky5g/ponzu/tokens"
	"mime/multipart"
)

type service struct {
	client driver.StorageClientInterface
	content.Service
}

type Service interface {
	content.Service
	GetAllUploads() ([]entities.FileUpload, error)
	GetFileUpload(target string) (*entities.FileUpload, error)
	DeleteFile(target string) error
	StoreFiles(files map[string]*multipart.FileHeader) (map[string]string, error)
	driver.StaticFileSystemInterface
}

func New(
	db driver.Database,
	searchClient driver.SearchClientInterface,
	client driver.StorageClientInterface) (Service, error) {
	uploadsRepository := db.Get(tokens.UploadRepositoryToken).(repositories.ContentRepositoryInterface)
	configRepository := db.Get(tokens.ConfigRepositoryToken).(repositories.ConfigRepositoryInterface)

	contentDomainService, err := content.New(uploadsRepository, configRepository, searchClient)
	if err != nil {
		return nil, err
	}

	if err = searchClient.CreateIndex(constants.UploadsEntityName, &entities.FileUpload{}); err != nil {
		return nil, err
	}

	s := &service{
		client:  client,
		Service: contentDomainService,
	}

	return s, nil
}
