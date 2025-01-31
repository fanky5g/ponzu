package request

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/util"
	"github.com/gorilla/schema"
)

var (
	PonzuRepeatPrefix          = ".__ponzu-repeat"
	PonzuFieldCollectionPrefix = ".__ponzu-field-collection"
)

func MapPayloadToGenericEntity(entity interface{}, payload map[string][]string) (interface{}, error) {
	addContentMetadata(payload)
	err := transformArrayFields(payload)
	if err != nil {
		return nil, err
	}

	dec := schema.NewDecoder()

	dec.SetAliasTag("json")
	dec.IgnoreUnknownKeys(true)
	err = dec.Decode(entity, payload)
	if err != nil {
		return nil, err
	}

	// we have to manually process field collections since gorilla/schema doesn't directly work with field collections
	if err = buildFieldCollections(entity, payload, dec); err != nil {
		return nil, err
	}

	cleanArrayFields(entity, payload)

	return entity, nil
}

func GetEntityFromFormData(entity interface{}, data url.Values) (interface{}, error) {
	return MapPayloadToGenericEntity(entity, data)
}

func GetEntity(entity interface{}, req *http.Request) (interface{}, error) {
	payload, err := GetRequestAsURLValues(req)
	if err != nil {
		return nil, err
	}

	return MapPayloadToGenericEntity(entity, payload)
}

func addContentMetadata(payload url.Values) {
	ts := payload.Get("timestamp")
	up := payload.Get("updated")

	// create a timestamp if one was not set
	if ts == "" {
		ts = fmt.Sprintf("%d", int64(time.Nanosecond)*time.Now().UTC().UnixNano()/int64(time.Millisecond))
		payload.Set("timestamp", ts)
	}

	if up == "" {
		payload.Set("updated", ts)
	}
}

func transformArrayFields(payload url.Values) error {
	// check for any multi-value fields (ex. checkbox fields)
	// and correctly format for storage. Essentially, we need
	// fieldX.0: value1, fieldX.1: value2 => fieldX: []string{value1, value2}
	fieldOrderValue := make(map[string]map[string][]string)
	for k, v := range payload {
		if strings.Contains(k, ".") {
			fo := strings.Split(k, ".")
			if len(fo) < 2 {
				continue
			}

			// put the order and the field value into map
			order := fo[len(fo)-1]
			field := strings.Join(fo[:len(fo)-1], ".")
			if _, err := strconv.ParseInt(order, 10, 64); err == nil {
				if len(fieldOrderValue[field]) == 0 {
					fieldOrderValue[field] = make(map[string][]string)
				}

				// orderValue is 0:[?type=Thing&id=1]
				orderValue := fieldOrderValue[field]
				orderValue[order] = v
				fieldOrderValue[field] = orderValue

				// discard the entity form value with name.N
				payload.Del(k)
			}
		}
	}

	// add/set the key & value to the entity form in order
	for f, ov := range fieldOrderValue {
		positions := make([]int, len(ov))
		i := 0
		for k := range ov {
			position, err := strconv.Atoi(k)
			if err != nil {
				// TODO: we should return error
				return fmt.Errorf("expected integer key in fieldOrderValue. Got %v", k)
			}

			positions[i] = position
			i = i + 1
		}

		sort.Ints(positions)

		for _, position := range positions {
			fieldValue := ov[fmt.Sprintf("%d", position)]

			if payload.Get(f) == "" {
				for j, fv := range fieldValue {
					if j == 0 {
						payload.Set(f, fv)
					} else {
						payload.Add(f, fv)
					}
				}
			} else {
				for _, fv := range fieldValue {
					payload.Add(f, fv)
				}
			}
		}
	}

	return nil
}

func cleanArrayFields(entity interface{}, payload url.Values) {
	repeatLengthIdentifier := make(map[string]int)
	repeatRemovedItemsIdentifier := make(map[string][]int)

	for k, v := range payload {
		repeatPrefix := ""
		if strings.HasPrefix(k, PonzuRepeatPrefix) {
			repeatPrefix = PonzuRepeatPrefix
		} else if strings.HasPrefix(k, PonzuFieldCollectionPrefix) {
			repeatPrefix = PonzuFieldCollectionPrefix
		}

		if repeatPrefix != "" {
			ponzuRepeatIdentifier := strings.TrimPrefix(k, repeatPrefix)

			if strings.HasSuffix(ponzuRepeatIdentifier, ".length") {
				if len(v) > 0 {
					if length, err := strconv.Atoi(v[0]); err == nil {
						mapKey := strings.TrimPrefix(
							strings.TrimSuffix(ponzuRepeatIdentifier, ".length"), ".",
						)
						repeatLengthIdentifier[mapKey] = length
					}
				}
			}

			if strings.HasSuffix(ponzuRepeatIdentifier, ".removed") {
				if len(v) > 0 {
					removedIndexesArray := strings.Split(v[0], ",")
					removedIndexesIntArray := make([]int, 0)
					for _, removedIndex := range removedIndexesArray {
						if index, err := strconv.Atoi(strings.TrimSpace(removedIndex)); err == nil {
							removedIndexesIntArray = append(removedIndexesIntArray, index)
						}
					}

					if len(removedIndexesIntArray) > 0 {
						mapKey := strings.TrimPrefix(
							strings.TrimSuffix(ponzuRepeatIdentifier, ".removed"), ".",
						)
						repeatRemovedItemsIdentifier[mapKey] = removedIndexesIntArray
					}
				}
			}

			payload.Del(k)
		}
	}

	for jsonFieldName, length := range repeatLengthIdentifier {
		fieldName := fieldNameByJSONTag(entity, jsonFieldName)
		if fieldName == "" {
			continue
		}

		v := reflect.Indirect(reflect.ValueOf(entity))
		field := v.FieldByName(fieldName)

		fieldCollections, isFieldCollections := (field.Interface()).(content.FieldCollections)
		if isFieldCollections {
			field = reflect.ValueOf(fieldCollections.Data())
		}

		if !field.IsValid() || util.SizeOfV(field) == length {
			continue
		}

		cleanedArray := reflect.MakeSlice(field.Type(), 0, length)
		if removedItems, ok := repeatRemovedItemsIdentifier[jsonFieldName]; ok {
			for i := 0; i < field.Len(); i++ {
				if !contains(removedItems, i) {
					cleanedArray = reflect.Append(cleanedArray, field.Index(i))
				}
			}
		}

		if isFieldCollections {
			fieldCollections.SetData(cleanedArray.Interface().([]content.FieldCollection))
		} else {
			field.Set(cleanedArray)
		}
	}
}

func buildFieldCollections(entity interface{}, payload map[string][]string, dec *schema.Decoder) error {
	v := reflect.ValueOf(entity)
	t := reflect.TypeOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		sField := t.Field(i)

		if field.IsValid() {
			fieldValue := field.Interface()
			fieldCollections, isFieldCollection := fieldValue.(content.FieldCollections)
			if !isFieldCollection {
				continue
			}

			jsonField, hasJsonField := sField.Tag.Lookup("json")
			if !hasJsonField {
				continue
			}

			allowedTypes := fieldCollections.AllowedTypes()
			for j, fieldCollection := range fieldCollections.Data() {
				valuePrefix := fmt.Sprintf("%s.%v.value", jsonField, j)
				values := make(map[string][]string)

				for entryKey, entryValue := range payload {
					if strings.HasPrefix(entryKey, valuePrefix) {
						key := strings.TrimPrefix(entryKey, valuePrefix)
						if strings.HasPrefix(key, ".") {
							key = strings.TrimPrefix(key, ".")
						}

						values[key] = entryValue
					}
				}

				if len(values) == 0 {
					continue
				}

				fieldCollection.Value = allowedTypes[fieldCollection.Type]()
				if err := dec.Decode(fieldCollection.Value, values); err != nil {
					return err
				}

				data := fieldCollections.Data()
				data[j] = fieldCollection
			}
		}
	}

	return nil
}

func fieldNameByJSONTag(p interface{}, jsonTagName string) string {
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag

		if jsonTag, ok := tag.Lookup("json"); ok {
			if jsonTag == jsonTagName {
				return typeField.Name
			}
		}
	}

	return ""
}

func contains(slice []int, val int) bool {
	for _, entry := range slice {
		if entry == val {
			return true
		}
	}
	return false
}
