package repository

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fanky5g/ponzu/database"
)

var testModelToken = "test"

type testEntity struct {
	ID      string    `json:"-"`
	Created time.Time `json:"-"`
	Updated time.Time `json:"-"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Age     int       `json:"age"`
}

func (entity *testEntity) EntityName() string {
	return "testEntity"
}

func (entity *testEntity) ItemID() string {
	return entity.ID
}

func (entity *testEntity) SetItemID(id string) {
	entity.ID = id
}

func (entity *testEntity) Title() string {
	return entity.Name
}

func (entity *testEntity) CreatedAt() int64 {
	if entity.Created.IsZero() {
		return 0
	}

	return entity.Created.Unix()
}

func (entity *testEntity) SetCreatedAt(t time.Time) {
	entity.Created = t
}

func (entity *testEntity) UpdatedAt() int64 {
	if entity.Updated.IsZero() {
		return 0
	}

	return entity.Updated.Unix()
}

func (entity *testEntity) SetUpdatedAt(t time.Time) {
	entity.Updated = t
}

type testModelDocument struct {
	*testEntity
}

func (document *testModelDocument) Value() (interface{}, error) {
	return json.Marshal(document)
}

func (document *testModelDocument) Scan(src interface{}) error {
	if byteSrc, ok := src.([]byte); ok {
		return json.Unmarshal(byteSrc, &document)
	}

	if stringSrc, ok := src.(string); ok {
		return json.NewDecoder(strings.NewReader(stringSrc)).Decode(&document)
	}

	return fmt.Errorf("unsupported type %T", src)
}

type testModel struct{}

func (*testModel) Name() string {
	return testModelToken
}

func (*testModel) NewEntity() interface{} {
	return new(testEntity)
}

func (model *testModel) ToDocument(entity interface{}) database.DocumentInterface {
	return &testModelDocument{
		testEntity: entity.(*testEntity),
	}
}
