package content

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
	"github.com/fanky5g/ponzu/tokens"
)

type service struct {
	repository   repositories.GenericRepositoryInterface
	searchClient driver.SearchClientInterface
	types        map[string]content.Builder
}

type Service interface {
	CreateContent(entityType string, entity interface{}) (string, error)
	DeleteContent(entityType, entityId string) error
	GetContent(entityType, entityId string) (interface{}, error)
	GetContentBySlug(slug string) (string, interface{}, error)
	GetAll(namespace string) ([]interface{}, error)
	GetAllWithOptions(namespace string, search *entities.Search) (int, []interface{}, error)
	UpdateContent(entityType, entityId string, update map[string]interface{}) (interface{}, error)
	ExportCSV(entityName string) (*entities.ResponseStream, error)
}

func New(
	db driver.Database,
	types map[string]content.Builder,
	searchClient driver.SearchClientInterface,
) (Service, error) {
	contentRepository := db.Get(tokens.ContentRepositoryToken).(repositories.GenericRepositoryInterface)

	for itemName, itemType := range types {
		if _, err := searchClient.GetIndex(itemName); err != nil {
			err = searchClient.CreateIndex(itemName, itemType())
			if err != nil {
				return nil, err
			}
		}
	}

	s := &service{
		repository:   contentRepository,
		searchClient: searchClient,
		types:        types,
	}

	return s, nil
}
