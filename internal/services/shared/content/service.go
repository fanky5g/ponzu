package content

import (
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
)

type service struct {
	repository       repositories.ContentRepositoryInterface
	configRepository repositories.ConfigRepositoryInterface
	searchClient     driver.SearchClientInterface
}

type Service interface {
	DeleteContent(entityType, entityId string) error
	CreateContent(entityType string, content interface{}) (string, error)
	UpdateContent(entityType, entityId string, update map[string]interface{}) (interface{}, error)
	GetContent(entityType, entityId string) (interface{}, error)
	GetContentBySlug(slug string) (string, interface{}, error)
	GetAllWithOptions(entityType string, search *entities.Search) (int, []interface{}, error)
	GetAll(entityType string) ([]interface{}, error)
}

func New(
	contentRepository repositories.ContentRepositoryInterface,
	configRepository repositories.ConfigRepositoryInterface,
	searchClient driver.SearchClientInterface,
) (Service, error) {
	s := &service{
		repository:       contentRepository,
		configRepository: configRepository,
		searchClient:     searchClient,
	}

	return s, nil
}
