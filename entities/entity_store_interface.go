package entities

import "github.com/fanky5g/ponzu/tokens"

type Persistable interface {
	GetRepositoryToken() tokens.RepositoryToken
}
