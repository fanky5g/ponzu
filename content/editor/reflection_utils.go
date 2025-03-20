package editor

import (
	"fmt"
	"github.com/fanky5g/ponzu/util"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func TagNameFromStructField(name string, post interface{}, args *FieldArgs) string {
	return tagNameFromStructField(name, post, args, 0)
}

func tagNameFromStructField(name string, post interface{}, args *FieldArgs, callDepth uint8) string {
	if name == "" {
		return name
	}

	if callDepth == 0 && args != nil && args.Parent != "" {
		name = strings.Join([]string{args.Parent, name}, ".")
	}

	parts := strings.Split(name, ".")
	fieldName := parts[0]
	v := reflect.ValueOf(post)
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()

		arrayEntityType := reflect.TypeOf(v.Interface()).Elem()
		if i, err := strconv.Atoi(fieldName); err == nil {
			if len(parts) > 1 {
				size := util.SizeOfV(v)
				var value interface{}
				if i < size {
					value = util.IndexAt(v, i)
				} else {
					value = util.MakeType(arrayEntityType)
				}

				nestedName := tagNameFromStructField(strings.Join(parts[1:], "."), value, args, callDepth+1)
				return strings.Join([]string{fieldName, nestedName}, ".")
			}

			return fieldName
		} else if isPositionalPlaceholder(fieldName, args) {
			if len(parts) > 1 {
				var value interface{}
				size := util.SizeOfV(v)
				if size > 0 {
					value = util.IndexAt(v, 0)
				} else {
					value = util.MakeType(t)
				}

				nestedName := tagNameFromStructField(strings.Join(parts[1:], "."), value, args, callDepth+1)
				return strings.Join([]string{fieldName, nestedName}, ".")
			}

			return fieldName
		}
	}

	field, ok := t.FieldByName(fieldName)
	if !ok {
		panic(
			"Couldn't get struct field for: " +
				fieldName +
				". Make sure you pass the right field name to editor field elements.",
		)
	}

	nestedName := ""
	if len(parts) > 1 {
		nestedName = tagNameFromStructField(
			strings.Join(parts[1:], "."),
			ValueByName(fieldName, post, nil).Interface(),
			args,
			callDepth+1,
		)
	}

	tag, ok := field.Tag.Lookup("json")
	if !ok {
		panic(
			"Couldn't get json struct tag for: " +
				name +
				". Struct fields for entities types must have 'json' tags.",
		)
	}

	if nestedName != "" {
		return strings.Join([]string{tag, nestedName}, ".")
	}

	return tag
}

// TagNameFromStructFieldMulti calls TagNameFromStructField and formats is for
// use with gorilla/schema
// due to the format in which gorilla/schema expects form names to be when
// one is associated with multiple values, we need to output the name as such.
// Ex. 'category.0', 'category.1', 'category.2' and so on.
func TagNameFromStructFieldMulti(name string, i int, post interface{}) string {
	tag := TagNameFromStructField(name, post, nil)

	return fmt.Sprintf("%s.%d", tag, i)
}

func ValueByName(name string, post interface{}, args *FieldArgs) reflect.Value {
	return valueByName(name, post, args, 0)
}

func valueByName(name string, post interface{}, args *FieldArgs, callDepth uint8) reflect.Value {
	if callDepth == 0 && args != nil && args.Parent != "" {
		name = strings.Join([]string{args.Parent, name}, ".")
	}

	parts := strings.Split(name, ".")
	fieldName := parts[0]
	v := reflect.ValueOf(post)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		if index, err := strconv.Atoi(fieldName); err == nil {
			v = getIndexAtPositionOrMakeNew(v, index)
		} else if isPositionalPlaceholder(fieldName, args) {
			v = getIndexAtPositionOrMakeNew(v, 0)
		}

		//else {
		//	v = v.Elem()
		//}

		if len(parts) > 1 {
			fieldName = parts[1]
			parts = parts[1:]
		} else {
			return v
		}
	}

	value := v.FieldByName(fieldName)
	if len(parts) > 1 {
		return valueByName(strings.Join(parts[1:], "."), value.Interface(), args, callDepth+1)
	}

	return value
}

func getIndexAtPositionOrMakeNew(v reflect.Value, idx int) reflect.Value {
	if util.SizeOfV(v) > 0 {
		return v.Index(idx)
	}

	return reflect.New(reflect.TypeOf(v.Interface()).Elem()).Elem()
}

func isPositionalPlaceholder(fieldName string, args *FieldArgs) bool {
	if args != nil {
		for _, positionalPlaceholder := range args.PositionalPlaceHolders {
			positionalPlaceholderRegexp := regexp.MustCompile(fmt.Sprintf("^%s$", positionalPlaceholder))
			if positionalPlaceholderRegexp.MatchString(fieldName) {
				return true
			}
		}
	}

	return false
}

// ValueFromStructField returns the value of a field in a struct
func ValueFromStructField(name string, post interface{}, args *FieldArgs) interface{} {
	field := ValueByName(name, post, args)

	switch field.Kind() {
	case reflect.String:
		return field.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%v", field.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%v", field.Uint())

	case reflect.Bool:
		return fmt.Sprintf("%t", field.Bool())

	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%v", field.Complex())

	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", field.Float())

	case reflect.Slice:
		t := reflect.TypeOf(field.Interface()).Elem()
		if t.Kind() == reflect.String {
			s := make([]string, field.Len())

			for i := 0; i < field.Len(); i++ {
				s[i] = field.Index(i).Interface().(string)
			}

			return s
		}

		return field.Interface()

	default:
		panic(fmt.Sprintf("Ponzu: Type '%s' for field '%s' not supported.", field.Type(), name))
	}
}

func makeEmptyType(p interface{}) interface{} {
	pType := reflect.TypeOf(p)
	if pType.Kind() == reflect.Pointer {
		pType = pType.Elem()
	}

	return reflect.New(pType).Interface()
}
