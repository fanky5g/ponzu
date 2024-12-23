package request

import (
	"net/url"

	"github.com/fanky5g/ponzu/entities"
)

func GetFileUploadFromFormData(data url.Values) (interface{}, error) {
	return MapPayloadToGenericEntity(new(entities.FileUpload), data)
}
