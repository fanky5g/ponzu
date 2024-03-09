package request

import (
	"github.com/fanky5g/ponzu/internal/domain/entities"
	"net/url"
)

func GetFileUploadFromFormData(data url.Values) (interface{}, error) {
	return mapPayloadToGenericEntity(func() interface{} {
		return new(entities.FileUpload)
	}, data)
}
