package content

import (
	"fmt"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/models"
	"github.com/fanky5g/ponzu/tokens"
	"log"
)

type service struct {
	repositories   map[string]driver.Repository
	slugRepository driver.Repository
	searchClient   driver.SearchClientInterface
	types          map[string]content.Builder
}

type Service interface {
	CreateContent(entityType string, entity interface{}) (string, error)
	DeleteContent(entityType, entityId string) error
	GetContent(entityType, entityId string) (interface{}, error)
	GetContentBySlug(slug string) (interface{}, error)
	GetAll(namespace string) ([]interface{}, error)
	GetAllWithOptions(namespace string, search *entities.Search) (int, []interface{}, error)
	UpdateContent(entityType, entityId string, update map[string]interface{}) (interface{}, error)
	ExportCSV(entityName string) (*entities.ResponseStream, error)
}

func (s *service) repository(entityType string) driver.Repository {
	repository := s.repositories[entityType]
	if repository == nil {
		log.Panicf("Failed to get repository for: %v", entityType)
	}

	return repository.(driver.Repository)
}

func New(
	db driver.Database,
	types map[string]content.Builder,
	searchClient driver.SearchClientInterface,
) (Service, error) {
	slugRepository := db.Get(
		models.WrapPonzuModelNameSpace(tokens.SlugRepositoryToken),
	).(driver.Repository)

	contentRepositories := make(map[string]driver.Repository)
	for itemName, itemType := range types {
		repository := db.Get(itemName)
		if repository == nil {
			return nil, fmt.Errorf("content repository for %s not implemented", itemName)
		}

		contentRepositories[itemName] = repository.(driver.Repository)

		if _, err := searchClient.GetIndex(itemName); err != nil {
			err = searchClient.CreateIndex(itemName, itemType())
			if err != nil {
				return nil, err
			}
		}
	}

	s := &service{
		repositories:   contentRepositories,
		slugRepository: slugRepository,
		searchClient:   searchClient,
		types:          types,
	}

	return s, nil
}
