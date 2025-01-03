package storage

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/pkg/errors"
)

// UploadFiles stores file uploads at paths like /YYYY/MM/filename.ext
func (s *service) UploadFiles(files map[string]*multipart.FileHeader) (map[string]string, error) {
	paths := make(map[string]string)
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

		paths[fileName] = u
		if err = s.storeFileInfo(size, fileName, u, fileHeader); err != nil {
			return nil, err
		}
	}

	return paths, nil
}

func (s *service) storeFileInfo(size int64, filename, urlPath string, file *multipart.FileHeader) error {
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
		return fmt.Errorf("error saving file storage record to database: %v", err)
	}

	if identifiable, ok := upload.(item.Identifiable); ok {
		if err = s.searchClient.Update(identifiable.ItemID(), upload); err != nil {
			return errors.Wrap(err, "Failed to update upload for search")
		}
	}

	return nil
}
