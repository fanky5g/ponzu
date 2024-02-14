package request

import (
	"fmt"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"net/http"
	"net/url"
)

func GetEntityFromFormData(entityType string, data url.Values) (interface{}, error) {
	// find the content type and decode values into it
	t, ok := item.Types[entityType]
	if !ok {
		return nil, fmt.Errorf(item.ErrTypeNotRegistered.Error(), entityType)
	}

	return mapPayloadToGenericEntity(t, data)
}

func GetEntity(entityType string, req *http.Request) (interface{}, error) {
	// find the content type and decode values into it
	t, ok := item.Types[entityType]
	if !ok {
		return nil, fmt.Errorf(item.ErrTypeNotRegistered.Error(), entityType)
	}

	payload, err := getRequestAsURLValues(req)
	if err != nil {
		return nil, err
	}

	return mapPayloadToGenericEntity(t, payload)
}
