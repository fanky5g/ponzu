package content

import (
	"log"

	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/search"
)

type Service struct {
	contentRepositories map[string]database.Repository
	slugs               database.Repository
	searchClient        search.Client
}

func (s *Service) repository(entityType string) database.Repository {
	repository, ok := s.contentRepositories[entityType]
	if !ok {
		log.Panicf("Failed to get repository for: %v", entityType)
	}

	return repository
}

func New(
	contentRepositories map[string]database.Repository,
	slugs database.Repository,
	searchClient search.Client,
) (*Service, error) {
	return &Service{
		contentRepositories: contentRepositories,
		slugs:               slugs,
		searchClient:        searchClient,
	}, nil
}
