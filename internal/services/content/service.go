package content

import (
	"github.com/fanky5g/ponzu/internal/domain/entities"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/domain/interfaces"
	"github.com/fanky5g/ponzu/internal/domain/services/content"
	"github.com/fanky5g/ponzu/internal/services"
)

var ServiceToken services.ServiceToken = "ContentService"

type service struct {
	contentDomainService content.Service
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
	content.Service
	ExportCSV(entityName string) (*entities.ResponseStream, error)
}

func New(
	contentRepository interfaces.ContentRepositoryInterface,
	configRepository interfaces.ConfigRepositoryInterface,
	searchClient interfaces.SearchClientInterface,
) (Service, error) {
	for itemName, itemType := range item.Types {
		if _, err := searchClient.GetIndex(itemName); err != nil {
			err = searchClient.CreateIndex(itemName, itemType())
			if err != nil {
				return nil, err
			}
		}
	}

	contentDomainService, err := content.New(contentRepository, configRepository, searchClient)
	if err != nil {
		return nil, err
	}

	s := &service{
		contentDomainService: contentDomainService,
	}

	return s, nil
}
