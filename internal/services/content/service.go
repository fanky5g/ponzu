package content

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
	contentService "github.com/fanky5g/ponzu/internal/services/shared/content"
	"github.com/fanky5g/ponzu/tokens"
)

type service struct {
	contentDomainService contentService.Service
	types                map[string]content.Builder
}

func (s *service) DeleteContent(entityType, entityId string) error {
	return s.contentDomainService.DeleteContent(entityType, entityId)
}

func (s *service) CreateContent(entityType string, content interface{}) (string, error) {
	return s.contentDomainService.CreateContent(entityType, content)
}

func (s *service) UpdateContent(entityType, entityId string, update map[string]interface{}) (interface{}, error) {
	return s.contentDomainService.UpdateContent(entityType, entityId, update)
}

func (s *service) GetContent(entityType, entityId string) (interface{}, error) {
	return s.contentDomainService.GetContent(entityType, entityId)
}

func (s *service) GetContentBySlug(slug string) (string, interface{}, error) {
	return s.contentDomainService.GetContentBySlug(slug)
}

func (s *service) GetAllWithOptions(entityType string, search *entities.Search) (int, []interface{}, error) {
	return s.contentDomainService.GetAllWithOptions(entityType, search)
}

func (s *service) GetAll(entityType string) ([]interface{}, error) {
	return s.contentDomainService.GetAll(entityType)
}

type Service interface {
	contentService.Service
	ExportCSV(entityName string) (*entities.ResponseStream, error)
}

func New(
	db driver.Database,
	types map[string]content.Builder,
	searchClient driver.SearchClientInterface,
) (Service, error) {
	contentRepository := db.Get(tokens.ContentRepositoryToken).(repositories.ContentRepositoryInterface)
	configRepository := db.Get(tokens.ConfigRepositoryToken).(repositories.ConfigRepositoryInterface)

	for itemName, itemType := range types {
		if _, err := searchClient.GetIndex(itemName); err != nil {
			err = searchClient.CreateIndex(itemName, itemType())
			if err != nil {
				return nil, err
			}
		}
	}

	contentDomainService, err := contentService.New(contentRepository, configRepository, searchClient)
	if err != nil {
		return nil, err
	}

	s := &service{
		contentDomainService: contentDomainService,
		types:                types,
	}

	return s, nil
}
