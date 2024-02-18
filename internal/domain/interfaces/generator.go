package interfaces

import (
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/domain/enum"
)

type ContentGenerator interface {
	Generate(contentType enum.ContentType, typeDefinition *item.TypeDefinition) error
	ValidateField(field *item.Field) error
}
