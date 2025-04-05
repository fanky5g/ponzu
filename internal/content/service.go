package content

import (
	"fmt"
	"log"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/fanky5g/ponzu/internal/content/dataexporter"
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/internal/datasource"
	"github.com/fanky5g/ponzu/internal/search"
	"github.com/fanky5g/ponzu/util"
	"github.com/pkg/errors"
)

type Service struct {
	repositories               map[string]database.Repository
	slugRepository             database.Repository
	searchClient               search.SearchInterface
	types                      map[string]content.Builder
	dataExporter               dataexporter.DataExporter
	uploadService              *UploadService
	workflowStateChangeHandler workflow.StateChangeTrigger
}

func New(
	db database.Database,
	types map[string]content.Builder,
	searchClient search.SearchInterface,
	dataExporter dataexporter.DataExporter,
	uploadService *UploadService,
	workflowStateChangeHandler workflow.StateChangeTrigger,
) (*Service, error) {
	slugRepository := db.GetRepositoryByToken(SlugRepositoryToken)

	contentRepositories := make(map[string]database.Repository)
	for entityName, ctor := range types {
		entity, ok := ctor().(database.Persistable)
		if !ok {
			return nil, fmt.Errorf("entity %s does not implement Persistable", entityName)
		}

		repository := db.GetRepositoryByToken(entity.GetRepositoryToken())
		if repository == nil {
			return nil, fmt.Errorf("content repository for %s not implemented", entityName)
		}

		contentRepositories[entityName] = repository
	}

	s := &Service{
		repositories:               contentRepositories,
		slugRepository:             slugRepository,
		searchClient:               searchClient,
		types:                      types,
		dataExporter:               dataExporter,
		uploadService:              uploadService,
		workflowStateChangeHandler: workflowStateChangeHandler,
	}

	return s, nil
}

func (s *Service) repository(entityType string) database.Repository {
	repository := s.repositories[entityType]
	if repository == nil {
		log.Panicf("Failed to get repository for: %v", entityType)
	}

	return repository
}

func (s *Service) ContentTypes() map[string]content.Builder {
	return s.types
}

func (s *Service) Type(name string) (interface{}, error) {
	ctor, ok := s.types[name]
	if !ok {
		return nil, ErrInvalidContentType
	}

	return ctor(), nil
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

	workflowStateManager, ok := entity.(workflow.LifecycleSupportedEntity)
	if ok && workflowStateManager.GetState() == "" {
		rootWorkflow, err := getContentWorkflow(entity)
		if err != nil {
			return "", err
		}

		workflowStateManager.SetState(rootWorkflow.GetState())
	}

	insert, err := repository.Insert(entity)
	if err != nil {
		return "", fmt.Errorf("failed to create content: %v", err)
	}

	identifiable = insert.(item.Identifiable)
	if _, err = s.slugRepository.Insert(&Slug{
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

func (s *Service) GetContentByIds(entityType string, entityIds ...string) ([]interface{}, error) {
	return s.repository(entityType).FindByIds(entityIds...)
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

	ss := match.(*Slug)
	return s.repository(ss.EntityType).FindOneById(ss.EntityId)
}

func (s *Service) GetSlug(entityType, entityId string) (*Slug, error) {
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

	return slug.(*Slug), nil
}

func (s *Service) GetAll(entityType string) ([]interface{}, error) {
	return s.repository(entityType).FindAll()
}

func (s *Service) GetAllWithOptions(entityType string, search *search.Search) ([]interface{}, int, error) {
	count, matches, err := s.repository(entityType).Find(search.SortOrder, search.Count, search.Offset)
	return matches, count, err
}

func (s *Service) GetNumberOfRows(entityType string) (int, error) {
	return s.repository(entityType).GetNumberOfRows()
}

func (s *Service) Export(entityName, exportType string) (datasource.Datasource, error) {
	t, ok := s.types[entityName]
	if !ok {
		return nil, fmt.Errorf(content.ErrTypeNotRegistered.Error(), entityName)
	}

	return s.dataExporter.Export(exportType, entityName, &contentDatasource{
		entity:         t(),
		contentService: s,
		entityName:     entityName,
	})
}

func (s *Service) ListReferences(typeName string, searchQuery *search.Search) ([]interface{}, int, error) {
	switch typeName {
	case UploadType:
		return s.uploadService.GetAllWithOptions(searchQuery)
	default:
		return s.GetAllWithOptions(typeName, searchQuery)
	}
}

func (s *Service) GetReference(entityName, entityId string) (interface{}, error) {
	switch entityName {
	case UploadType:
		return s.uploadService.GetUpload(entityId)
	default:
		return s.GetContent(entityName, entityId)
	}
}

func (s *Service) GetReferences(entityName string, entityIds ...string) ([]interface{}, error) {
	switch entityName {
	case UploadType:
		ups, err := s.uploadService.GetUploads(entityIds...)
		if err != nil {
			return nil, err
		}

		out := make([]interface{}, len(ups))
		for i, up := range ups {
			out[i] = up
		}

		return out, nil
	default:
		return s.GetContentByIds(entityName, entityIds...)
	}
}
