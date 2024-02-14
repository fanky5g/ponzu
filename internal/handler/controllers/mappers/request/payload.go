package request

import (
	"net/http"
	"net/url"
)

func getRequestAsURLValues(req *http.Request) (url.Values, error) {
	payload := make(url.Values)
	var err error

	switch getContentType(req) {
	case "application/x-www-form-urlencoded":
		payload = req.URL.Query()
		break
	case "multipart/form-data":
		if err = req.ParseMultipartForm(1024 * 1024 * 4); err != nil {
			return nil, err
		}

		payload = req.PostForm
		break
	case "application/json":
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
