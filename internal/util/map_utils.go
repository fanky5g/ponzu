package util

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
)

func JSONMapToURLValues(payload map[string]interface{}) url.Values {
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

			break
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

	if data.Get(k) == "" {
		data.Set(k, str)
	} else {
		data.Add(k, str)
	}
}
