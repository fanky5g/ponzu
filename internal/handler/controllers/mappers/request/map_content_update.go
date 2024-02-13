package request

import (
	"net/http"
)

func MapRequestToContentUpdate(req *http.Request) (map[string]interface{}, error) {
	payload, err := getRequestPayload(req)
	if err != nil {
		return nil, err
	}

	update := make(map[string]interface{})
	for k, v := range payload {
		update[k] = v
	}

	return update, nil
}
