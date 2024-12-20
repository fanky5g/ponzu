package content

import (
	"fmt"
	"log"
	"time"

	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/datasource"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/fanky5g/ponzu/util"
	"github.com/pkg/errors"
)

var CSVChunkSize = 50

type Service struct {
	repositories                     map[string]driver.Repository
	slugRepository                   driver.Repository
	searchClient                     driver.SearchInterface
	types                            map[string]content.Builder
	csvExportDatasourceReaderFactory datasource.DataSourceReaderFactory
}

type CSVExportDataSource struct {
	entity         interface{}
	entityName     string
	contentService *Service
}

func New(
	db driver.Database,
	csvDataSourceFactory datasource.DataSourceReaderFactory,
	types map[string]content.Builder,
	searchClient driver.SearchInterface,
) (*Service, error) {
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

	s := &Service{
		repositories:   contentRepositories,
		slugRepository: slugRepository,
		searchClient:   searchClient,
		types:          types,
	}

	return s, nil
}

func (s *Service) repository(entityType string) driver.Repository {
	repository := s.repositories[entityType]
	if repository == nil {
		log.Panicf("Failed to get repository for: %v", entityType)
	}

	return repository
}

func (s *Service) CreateContent(entityType string, entity interface{}) (string, error) {
	repository := s.repository(entityType)
	identifiable, ok := entity.(item.Identifiable)
	if !ok {
		return "", errors.New("item does not implement identifiable interface")
	}

	sluggable, ok := entity.(item.Sluggable)
	if !ok {
		return "", errors.New("entity does not implement sluggable interface")
	}

	if sluggable.ItemSlug() == "" {
		slug, err := util.Slugify(sluggable.GetTitle())
		if err != nil {
			return "", fmt.Errorf("failed to get slug: %v", err)
		}

		sluggable.SetSlug(slug)
	}

	if workflowStateManager, ok := entity.(workflow.StateManager); ok {
		workflowStateManager.SetState(workflow.DraftState)
	}

	content, err := repository.Insert(entity)
	if err != nil {
		return "", fmt.Errorf("failed to create content: %v", err)
	}

	identifiable = content.(item.Identifiable)
	if _, err = s.slugRepository.Insert(&entities.Slug{
		EntityType: entityType,
		EntityId:   identifiable.ItemID(),
		Slug:       sluggable.ItemSlug(),
	}); err != nil {
		return "", fmt.Errorf("failed to save slug: %v", err)
	}

	if err = s.searchClient.Update(identifiable.ItemID(), entity); err != nil {
		return "", fmt.Errorf("failed to index entity: %v", err)
	}

	return identifiable.ItemID(), nil
}

// UpdateContent supports only full updates. The entire structure will be overwritten.
func (s *Service) UpdateContent(entityType, entityId string, update interface{}) (interface{}, error) {
	identifiable, ok := update.(item.Identifiable)
	if !ok {
		return nil, fmt.Errorf("update not supported for %s", entityType)
	}

	identifiable.SetItemID(entityId)
	var sluggable item.Sluggable
	if sluggable, ok = update.(item.Sluggable); ok {
		slug, err := s.GetSlug(entityType, entityId)
		if err != nil {
			return nil, err
		}

		if slug != nil {
			sluggable.SetSlug(slug.Slug)
		}
	}

	update, err := s.repository(entityType).UpdateById(entityId, update)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to update content")
	}

	if err = s.searchClient.Update(entityId, update); err != nil {
		return nil, errors.Wrap(err, "Failed to update search index")
	}

	return update, nil
}

func (s *Service) DeleteContent(entityType string, entityIds ...string) error {
	repository := s.repository(entityType)

	for _, entityId := range entityIds {
		if err := repository.DeleteById(entityId); err != nil {
			return fmt.Errorf("failed to delete: %v", err)
		}

		if err := s.slugRepository.DeleteBy("entity_id", constants.Equal, entityId); err != nil {
			return fmt.Errorf("failed to delete slug: %v", err)
		}

		if err := s.searchClient.Delete(entityType, entityId); err != nil {
			return errors.Wrap(err, "failed to delete indexed content")
		}

	}

	return nil
}

func (s *Service) GetContent(entityType, entityId string) (interface{}, error) {
	return s.repository(entityType).FindOneById(entityId)
}

func (s *Service) GetContentBySlug(slug string) (interface{}, error) {
	match, err := s.slugRepository.FindOneBy(map[string]interface{}{
		"slug": slug,
	})

	if err != nil {
		return nil, err
	}

	if match == nil {
		return nil, nil
	}

	ss := match.(*entities.Slug)
	return s.repository(ss.EntityType).FindOneById(ss.EntityId)
}

func (s *Service) GetSlug(entityType, entityId string) (*entities.Slug, error) {
	slug, err := s.slugRepository.FindOneBy(map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityId,
	})

	if err != nil {
		return nil, err
	}

	if slug == nil {
		return nil, nil
	}

	return slug.(*entities.Slug), nil
}

func (s *Service) GetAll(entityType string) ([]interface{}, error) {
	return s.repository(entityType).FindAll()
}

func (s *Service) GetAllWithOptions(entityType string, search *entities.Search) ([]interface{}, int, error) {
	count, matches, err := s.repository(entityType).Find(search.SortOrder, search.Pagination)
	return matches, count, err
}

func (s *Service) GetNumberOfRows(entityType string) (int, error) {
	return s.repository(entityType).GetNumberOfRows()
}

func (s *Service) ExportCSV(entityName string) (*entities.ResponseStream, error) {
	t, ok := s.types[entityName]
	if !ok {
		return nil, fmt.Errorf(content.ErrTypeNotRegistered.Error(), entityName)
	}

	entity := t()
	csvExportDataSource, err := NewCSVExportDataSource(entity)
	if err != nil {
		return nil, err
	}

	r, err := s.csvExportDatasourceReaderFactory(csvExportDataSource)
	if err != nil {
		return nil, err
	}

	return &entities.ResponseStream{
		ContentType:        "text/csv",
		ContentDisposition: fmt.Sprintf(`attachment; filename="export-%s-%d.csv"`, entityName, time.Now().Unix()),
		Payload:            r,
	}, nil
}
