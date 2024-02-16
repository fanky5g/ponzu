package interfaces

import (
	"github.com/fanky5g/ponzu/internal/domain/entities"
	"github.com/fanky5g/ponzu/internal/domain/enum"
)

type ContentGenerator interface {
	Generate(contentType enum.ContentType, typeDefinition *entities.TypeDefinition) error
	ValidateField(field *entities.Field) error
}
