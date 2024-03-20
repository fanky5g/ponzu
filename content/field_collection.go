package content

type FieldCollections interface {
	Name() string
	AllowedTypes() map[string]Builder
	Data() []FieldCollection
	Add(fieldCollection FieldCollection)
	Set(i int, fieldCollection FieldCollection)
	SetData(data []FieldCollection)
}

type FieldCollection struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
