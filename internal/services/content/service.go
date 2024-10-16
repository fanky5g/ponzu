package content

import (
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/tokens"
	"log"
)

type service struct {
	repositories   map[string]driver.Repository
	slugRepository driver.Repository
	searchClient   driver.SearchInterface
	types          map[string]content.Builder
}

type Service interface {
	CreateContent(entityType string, entity interface{}) (string, error)
	DeleteContent(entityType string, entityIds ...string) error
	GetContent(entityType, entityId string) (interface{}, error)
	GetContentBySlug(slug string) (interface{}, error)
	GetAll(namespace string) ([]interface{}, error)
	GetAllWithOptions(namespace string, search *entities.Search) ([]interface{}, int, error)
	UpdateContent(entityType, entityId string, update interface{}) (interface{}, error)
	ExportCSV(entityName string) (*entities.ResponseStream, error)
}

func (s *service) repository(entityType string) driver.Repository {
	repository := s.repositories[entityType]
	if repository == nil {
		log.Panicf("Failed to get repository for: %v", entityType)
	}

	return repository
}

func New(
	db driver.Database,
	types map[string]content.Builder,
	searchClient driver.SearchInterface,
) (Service, error) {
	slugRepository := db.GetRepositoryByToken(tokens.SlugRepositoryToken)

	contentRepositories := make(map[string]driver.Repository)
	for entityName, entityConstructor := range types {
		entity := entityConstructor()
		persistable, ok := entity.(entities.Persistable)
		if !ok {
			return nil, fmt.Errorf("entity %s does not implement Persistable", entityName)
		}

		repository := db.GetRepositoryByToken(persistable.GetRepositoryToken())
		if repository == nil {
			return nil, fmt.Errorf("content repository for %s not implemented", entityName)
		}

		contentRepositories[entityName] = repository
	}

	s := &service{
		repositories:   contentRepositories,
		slugRepository: slugRepository,
		searchClient:   searchClient,
		types:          types,
	}

	return s, nil
}
