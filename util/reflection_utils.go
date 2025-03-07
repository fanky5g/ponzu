package util

import (
	"fmt"
	"reflect"
)

func SizeOfV(v reflect.Value) int {
	if v.IsZero() || v.IsNil() {
		return 0
	}

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if !v.IsValid() {
		panic(fmt.Sprintf("cannot compute size of invalid value: %+v", v))
	}

	return v.Len()
}

func MakeType(pType reflect.Type) interface{} {
	if pType.Kind() == reflect.Pointer {
		pType = pType.Elem()
	}

	return reflect.New(pType).Interface()
}

func IndexAt(v reflect.Value, i int) interface{} {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	return v.Index(i).Interface()
}

func FieldByJSONTagName(structType interface{}, jsonTagName string) reflect.Value {
	v := reflect.ValueOf(structType)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag

		if jsonTag, ok := tag.Lookup("json"); ok {
			if jsonTag == jsonTagName {
				return v.FieldByName(typeField.Name)
			}
		}
	}

	return reflect.Value{}
}

func StringFieldByJSONTagName(structType interface{}, jsonTagName string) (string, error) {
	if structType == nil {
		return "", nil
	}

	value := FieldByJSONTagName(structType, jsonTagName)
	if !value.IsValid() {
		return "", fmt.Errorf("%s is not a valid field", jsonTagName)
	}

	if value.Kind() != reflect.String {
		return "", fmt.Errorf("%s is not a string", jsonTagName)
	}

	return value.String(), nil
}
