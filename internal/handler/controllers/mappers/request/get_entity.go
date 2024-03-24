package request

import (
	"github.com/fanky5g/ponzu/content"
	"net/http"
	"net/url"
)

func GetEntityFromFormData(t content.Builder, data url.Values) (interface{}, error) {
	return mapPayloadToGenericEntity(t, data)
}

func GetEntity(t content.Builder, req *http.Request) (interface{}, error) {
	payload, err := getRequestAsURLValues(req)
	if err != nil {
		return nil, err
	}

	return mapPayloadToGenericEntity(t, payload)
}
