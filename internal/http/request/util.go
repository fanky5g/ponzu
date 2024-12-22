package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)


var ErrUnsupportedContentType  = errors.New("unsupported content type")

func GetContentType(req *http.Request) string {
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, ";") {
		contentType = strings.TrimSpace(contentType[:strings.Index(contentType, ";")])
	}

	return contentType
}

func GetRequestAsURLValues(req *http.Request) (url.Values, error) {
	payload := make(url.Values)
	var err error

	switch GetContentType(req) {
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

func mapJSONContentToURLValues(req *http.Request) (map[string][]string, error) {
	var payload map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return MapToURLValues(payload), nil
}

func MapToURLValues(payload map[string]interface{}) url.Values {
	if payload == nil || len(payload) == 0 {
		return nil
	}

	data := make(url.Values)
	for key, value := range payload {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Array:
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				setInterface(key, v.Index(i).Interface(), data)
			}
		default:
			setInterface(key, value, data)
		}
	}

	return data
}

func setInterface(k string, v interface{}, data url.Values) {
	var str string
	kind := reflect.ValueOf(v).Kind()
	switch kind {
	case reflect.String:
		str = v.(string)
	case reflect.Bool:
		str = fmt.Sprint(v.(bool))
	case reflect.Float64:
		str = fmt.Sprint(v.(float64))
	case reflect.Map:
		for key, value := range v.(map[string]interface{}) {
			setInterface(fmt.Sprintf("%s.%s", k, key), value, data)
		}
		return
	default:
		log.Println("Unsupported field", k, kind)
	}

	if data.Get(k) == "" {
		data.Set(k, str)
	} else {
		data.Add(k, str)
	}
}
