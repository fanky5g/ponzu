package content

const SlugRepositoryToken = "slugs"

type Slug struct {
	EntityType string `json:"entity_type"`
	EntityId   string `json:"entity_id"`
	Slug       string `json:"slug"`
}

func (*Slug) GetRepositoryToken() string {
	return SlugRepositoryToken
}

func (*Slug) EntityName() string {
	return "Slug"
}
