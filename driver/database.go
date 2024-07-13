package driver

import "github.com/fanky5g/ponzu/tokens"

type Database interface {
	GetRepositoryByToken(token tokens.RepositoryToken) Repository
	Close() error
}
