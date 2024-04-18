package entities

import "github.com/fanky5g/ponzu/tokens"

type EntityStoreInterface interface {
	GetRepositoryToken() tokens.RepositoryToken
}
