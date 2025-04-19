package content

import (
	"fmt"
	"github.com/fanky5g/ponzu/internal/storage"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/fanky5g/ponzu/content/entities"
	contentEntities "github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/internal/search"
	"github.com/pkg/errors"
)

var UploadType = "Upload"

type UploadService struct {
	client       storage.Client
	repository   database.Repository
	searchClient search.SearchInterface
}

func NewUploadService(
	db database.Database,
	searchClient search.SearchInterface,
	client storage.Client) (*UploadService, error) {
	s := &UploadService{
		client:       client,
		searchClient: searchClient,
		repository:   db.GetRepositoryByToken(contentEntities.UploadRepositoryToken),
	}

	return s, nil
}

func (s *UploadService) GetUpload(entityId string) (*contentEntities.Upload, error) {
	file, err := s.repository.FindOneById(entityId)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, nil
	}

	return file.(*contentEntities.Upload), nil
}

func (s *UploadService) GetUploads(entityIds ...string) ([]*contentEntities.Upload, error) {
	matches, err := s.repository.FindByIds(entityIds...)
	if err != nil {
		return nil, err
	}

	out := make([]*contentEntities.Upload, len(matches))
	for i, match := range matches {
		out[i] = match.(*contentEntities.Upload)
	}

	return out, nil
}

func (s *UploadService) GetAllWithOptions(search *search.Search) ([]interface{}, int, error) {
	total, files, err := s.repository.Find(search.SortOrder, search.Count, search.Offset)
	if err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

func (s *UploadService) Open(name string) (http.File, error) {
	return s.client.Open(name)
}

func (s *UploadService) DeleteUpload(entityIds ...string) error {
	for _, entityId := range entityIds {
		f, err := s.GetUpload(entityId)
		if err != nil {
			return err
		}

		if f == nil {
			return nil
		}

		if err = s.client.Delete(f.Path); err != nil {
			return fmt.Errorf("failed to delete from file store: %v", err)
		}

		if err = s.repository.DeleteById(entityId); err != nil {
			return errors.Wrap(err, "Failed to delete item from database")
		}

		if err = s.searchClient.Delete(constants.UploadEntityName, entityId); err != nil {
			return errors.Wrap(err, "Failed to delete search index entry")
		}
	}

	return nil
}

// UploadFiles stores file uploads at paths like /YYYY/MM/filename.ext
func (s *UploadService) UploadFiles(files map[string]*multipart.FileHeader) ([]*entities.Upload, error) {
	uploadedFiles := make([]*entities.Upload, len(files))
	i := 0
	for name, fileHeader := range files {
		nameParts := strings.Split(name, ":")
		fileName := nameParts[0]
		if len(nameParts) > 1 {
			fileName = nameParts[1]
		}

		f, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("could not open file for uploading: %v", err)
		}

		u, size, err := s.client.Save(fileName, f)
		if err != nil {
			return nil, err
		}

		var upload *entities.Upload
		upload, err = s.storeFileInfo(size, fileName, u, fileHeader)
		if err != nil {
			return nil, err
		}

		uploadedFiles[i] = upload
		i = i + 1
	}

	return uploadedFiles, nil
}

func (s *UploadService) storeFileInfo(size int64, filename, urlPath string, file *multipart.FileHeader) (*entities.Upload, error) {
	ts := int64(time.Nanosecond) * time.Now().UTC().UnixNano() / int64(time.Millisecond)
	entity := &entities.Upload{
		Name:          filename,
		Path:          urlPath,
		ContentLength: size,
		ContentType:   file.Header.Get("Content-Type"),
		Item: item.Item{
			Timestamp: ts,
			Updated:   ts,
		},
	}

	upload, err := s.repository.Insert(entity)
	if err != nil {
		return nil, fmt.Errorf("error saving file storage record to database: %v", err)
	}

	if identifiable, ok := upload.(item.Identifiable); ok {
		if err = s.searchClient.Update(identifiable.ItemID(), upload); err != nil {
			return nil, errors.Wrap(err, "Failed to update upload for search")
		}
	}

	return entity, nil
}
