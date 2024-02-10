package interfaces

import (
	"github.com/fanky5g/ponzu/internal/domain/entities"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/domain/enum"
)

type CRUDInterface interface {
	SetEntity(entityType string, entity interface{}) (string, error)
	DeleteEntity(entityId string) error
	FindByTarget(targets []string) ([]interface{}, error)
	FindOneByTarget(target string) (interface{}, error)
	FindOneBySlug(slug string) (string, interface{}, error)
	FindAll(namespace string) ([]interface{}, error)
	FindAllWithOptions(
		namespace string,
		order enum.SortOrder,
		pagination *entities.Pagination,
	) (int, []interface{}, error)
}

type EntityIdentifierInterface interface {
	UniqueSlug(slug string) (string, error)
	IsValidID(id string) bool
	NextIDSequence(entityType string) (string, error)
}

type ContentRepositoryInterface interface {
	CreateEntityStore(entityName string, entityType interface{}) error
	CRUDInterface
	EntityIdentifierInterface
	Types() map[string]item.EntityBuilder
}
