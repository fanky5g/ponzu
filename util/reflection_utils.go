package util

import (
	"fmt"
	"reflect"
)

func SizeOfV(v reflect.Value) int {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if !v.IsValid() {
		panic(fmt.Sprintf("cannot compute size of invalid value %v", v))
	}

	if v.IsNil() || v.IsZero() {
		return 0
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
