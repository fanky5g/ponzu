package models

import (
	"github.com/google/uuid"
	"time"
)

type Model struct {
	ID        uuid.UUID         `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt time.Time         `json:"deleted_at"`
	Document  DocumentInterface `json:"document"`
}

type ModelInterface interface {
	Name() string
	ToDocument(entity interface{}) DocumentInterface
	NewEntity() interface{}
}
