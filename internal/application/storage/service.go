// Package storage provides a re-usable file storage and storage utility for Ponzu
// systems to handle multipart form data.
package storage

import (
	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/search"
	"github.com/fanky5g/ponzu/storage"
)

type Service struct {
	client       storage.Client
	uploads      database.Repository
	searchClient search.Client
}

func New(uploads database.Repository, searchClient search.Client, storageClient storage.Client) (*Service, error) {
	return &Service{
		client:       storageClient,
		searchClient: searchClient,
		uploads:      uploads,
	}, nil
}
