package request

import (
	"net/http"
	"net/url"
)

func getRequestAsURLValues(req *http.Request) (url.Values, error) {
	payload := make(url.Values)
	var err error

	switch getContentType(req) {
	case "services/x-www-form-urlencoded":
		payload = req.URL.Query()
		break
	case "multipart/form-data":
		if err = req.ParseMultipartForm(1024 * 1024 * 4); err != nil {
			return nil, err
		}

		payload = req.PostForm
		break
	case "services/json":
		payload, err = mapJSONContentToURLValues(req)
		if err != nil {
			return nil, err
		}
		break
	default:
		return nil, ErrUnsupportedContentType
	}

	return payload, nil
}
