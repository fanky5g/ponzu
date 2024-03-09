package root

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/domain/interfaces"
)

var contentIndexName = "__contentIndex"

type repository struct {
	db        *bolt.DB
	entityMap map[string]item.EntityBuilder
}

// New instantiates common repository functions implemented by all repositories
func New(db *bolt.DB, entityMap map[string]item.EntityBuilder) (interfaces.ContentRepositoryInterface, error) {
	repo := &repository{db: db, entityMap: entityMap}
	for itemName, itemType := range entityMap {
		if err := repo.CreateEntityStore(itemName, itemType()); err != nil {
			return nil, err
		}
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(contentIndexName))
		return err
	}); err != nil {
		return nil, fmt.Errorf("failed to create storage bucket: %v", contentIndexName)
	}

	return repo, nil
}
