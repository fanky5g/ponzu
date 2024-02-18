package util

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
)

// TODO: support for nested arrays
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
