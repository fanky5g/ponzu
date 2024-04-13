package entities

type Slug struct {
	EntityType string `json:"entity_type"`
	EntityId   string `json:"entity_id"`
	Slug       string `json:"slug"`
}
