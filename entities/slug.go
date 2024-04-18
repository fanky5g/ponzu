package entities

import "github.com/fanky5g/ponzu/tokens"

type Slug struct {
	EntityType string `json:"entity_type"`
	EntityId   string `json:"entity_id"`
	Slug       string `json:"slug"`
}

func (*Slug) GetRepositoryToken() tokens.RepositoryToken {
	return tokens.SlugRepositoryToken
}
