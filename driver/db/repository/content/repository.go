package content

import (
	"github.com/boltdb/bolt"
	"github.com/fanky5g/ponzu/driver/db/repository/root"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/domain/interfaces"
)

func New(db *bolt.DB) (interfaces.ContentRepositoryInterface, error) {
	return root.New(db, item.Types)
}
