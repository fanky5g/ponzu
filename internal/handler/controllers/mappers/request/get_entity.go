package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

var ErrUnsupportedContentType = errors.New("unsupported content type")

func GetEntityFromFormData(entityType string, data url.Values) (interface{}, error) {
	// find the content type and decode values into it
	t, ok := item.Types[entityType]
	if !ok {
		return nil, fmt.Errorf(item.ErrTypeNotRegistered.Error(), entityType)
	}

	return getGenericEntity(t, data)
}

func GetEntity(entityType string, req *http.Request) (interface{}, error) {
	// find the content type and decode values into it
	t, ok := item.Types[entityType]
	if !ok {
		return nil, fmt.Errorf(item.ErrTypeNotRegistered.Error(), entityType)
	}

	contentType := getContentType(req)

	if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
		if err := req.ParseMultipartForm(1024 * 1024 * 4); err != nil {
			return nil, err
		}

		return getGenericEntity(t, req.PostForm)
	} else if contentType == "application/json" {
		c, err := mapJSONContentToURLValues(req)
		if err != nil {
			return nil, err
		}

		return getGenericEntity(t, c)
	}

	return nil, ErrUnsupportedContentType
}

func get(key string, values map[string][]string) string {
	vs := values[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func set(key, value string, values map[string][]string) {
	values[key] = []string{value}
}

func add(key, value string, values map[string][]string) {
	values[key] = append(values[key], value)
}

func del(key string, values map[string][]string) {
	delete(values, key)
}

func getGenericEntity(t item.EntityBuilder, data map[string][]string) (interface{}, error) {
	entity := t()
	ts := get("timestamp", data)
	up := get("updated", data)

	// create a timestamp if one was not set
	if ts == "" {
		ts = fmt.Sprintf("%d", int64(time.Nanosecond)*time.Now().UTC().UnixNano()/int64(time.Millisecond))
		set("timestamp", ts, data)
	}

	if up == "" {
		set("updated", ts, data)
	}

	// check for any multi-value fields (ex. checkbox fields)
	// and correctly format for storage. Essentially, we need
	// fieldX.0: value1, fieldX.1: value2 => fieldX: []string{value1, value2}
	fieldOrderValue := make(map[string]map[string][]string)
	for k, v := range data {
		if strings.Contains(k, ".") {
			fo := strings.Split(k, ".")

			// put the order and the field value into map
			field := fo[0]
			order := fo[1]
			if len(fieldOrderValue[field]) == 0 {
				fieldOrderValue[field] = make(map[string][]string)
			}

			// orderValue is 0:[?type=Thing&id=1]
			orderValue := fieldOrderValue[field]
			orderValue[order] = v
			fieldOrderValue[field] = orderValue

			// discard the entity form value with name.N
			del(k, data)
		}
	}

	// add/set the key & value to the entity form in order
	for f, ov := range fieldOrderValue {
		for i := 0; i < len(ov); i++ {
			position := fmt.Sprintf("%d", i)
			fieldValue := ov[position]

			if get(f, data) == "" {
				for i, fv := range fieldValue {
					if i == 0 {
						set(f, fv, data)
					} else {
						add(f, fv, data)
					}
				}
			} else {
				for _, fv := range fieldValue {
					add(f, fv, data)
				}
			}
		}
	}

	dec := schema.NewDecoder()
	dec.SetAliasTag("json")     // allows simpler struct tagging when creating a content type
	dec.IgnoreUnknownKeys(true) // will skip over form values submitted, but not in struct
	err := dec.Decode(entity, data)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func mapJSONContentToURLValues(req *http.Request) (map[string][]string, error) {
	data := make(map[string][]string)
	var requestBody map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		return nil, err
	}

	for key, value := range requestBody {
		switch reflect.ValueOf(value).Kind() {
		case reflect.Array:
		case reflect.Slice:
			arr := value.([]interface{})
			for _, arrValue := range arr {
				setInterface(key, arrValue, data)
			}
			break
		default:
			setInterface(key, value, data)
		}
	}

	return data, nil
}

func setInterface(k string, v interface{}, data map[string][]string) {
	var str string
	kind := reflect.ValueOf(v).Kind()
	switch kind {
	case reflect.String:
		str = v.(string)
		break
	case reflect.Bool:
		str = fmt.Sprint(v.(bool))
		break
	case reflect.Float64:
		str = fmt.Sprint(v.(float64))
		break
	default:
		log.Println("Unsupported field", k, kind)
	}

	if get(k, data) == "" {
		set(k, str, data)
	} else {
		add(k, str, data)
	}
}
