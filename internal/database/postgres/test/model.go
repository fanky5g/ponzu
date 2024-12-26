package test

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fanky5g/ponzu/database"
)

var ModelToken = "test"

type Entity struct {
	ID      string    `json:"-"`
	Created time.Time `json:"-"`
	Updated time.Time `json:"-"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Age     int       `json:"age"`
}

func (entity *Entity) EntityName() string {
	return "Entity"
}

func (entity *Entity) ItemID() string {
	return entity.ID
}

func (entity *Entity) SetItemID(id string) {
	entity.ID = id
}

func (entity *Entity) Title() string {
	return entity.Name
}

func (entity *Entity) CreatedAt() int64 {
	if entity.Created.IsZero() {
		return 0
	}

	return entity.Created.Unix()
}

func (entity *Entity) SetCreatedAt(t time.Time) {
	entity.Created = t
}

func (entity *Entity) UpdatedAt() int64 {
	if entity.Updated.IsZero() {
		return 0
	}

	return entity.Updated.Unix()
}

func (entity *Entity) SetUpdatedAt(t time.Time) {
	entity.Updated = t
}

func (entity *Entity) GetRepositoryToken() string {
	return "entity"
}

type ModelDocument struct {
	*Entity
}

func (document *ModelDocument) Value() (interface{}, error) {
	return json.Marshal(document)
}

func (document *ModelDocument) Scan(src interface{}) error {
	if byteSrc, ok := src.([]byte); ok {
		return json.Unmarshal(byteSrc, &document)
	}

	if stringSrc, ok := src.(string); ok {
		return json.NewDecoder(strings.NewReader(stringSrc)).Decode(&document)
	}

	return fmt.Errorf("unsupported type %T", src)
}

type Model struct{}

func (*Model) Name() string {
	return ModelToken
}

func (*Model) NewEntity() interface{} {
	return new(Entity)
}

func (model *Model) ToDocument(entity interface{}) database.DocumentInterface {
	return &ModelDocument{
		Entity: entity.(*Entity),
	}
}
