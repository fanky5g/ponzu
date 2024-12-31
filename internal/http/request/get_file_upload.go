package request

import (
	"net/url"

	"github.com/fanky5g/ponzu/content/entities"
)

func GetFileUploadFromFormData(data url.Values) (interface{}, error) {
	return MapPayloadToGenericEntity(new(entities.FileUpload), data)
}
