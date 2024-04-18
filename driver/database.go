package driver

import (
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/tokens"
)

type Database interface {
	GetRepositoryByToken(token tokens.RepositoryToken) Repository
	GetRepository(entity entities.EntityStoreInterface) Repository
	Close() error
}
