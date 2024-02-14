package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/util"
	"github.com/gorilla/schema"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrUnsupportedContentType = errors.New("unsupported content type")

func mapPayloadToGenericEntity(t item.EntityBuilder, payload map[string][]string) (interface{}, error) {
	entity := t()
	addContentMetadata(payload)
	applyContentFieldTransforms(payload)

	dec := schema.NewDecoder()
	dec.SetAliasTag("json")     // allows simpler struct tagging when creating a content type
	dec.IgnoreUnknownKeys(true) // will skip over form values submitted, but not in struct
	err := dec.Decode(entity, payload)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func mapJSONContentToURLValues(req *http.Request) (map[string][]string, error) {
	var payload map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return util.JSONMapToURLValues(payload), nil
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

func applyContentFieldTransforms(payload url.Values) {
	// check for any multi-value fields (ex. checkbox fields)
	// and correctly format for storage. Essentially, we need
	// fieldX.0: value1, fieldX.1: value2 => fieldX: []string{value1, value2}
	fieldOrderValue := make(map[string]map[string][]string)
	for k, v := range payload {
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
			payload.Del(k)
		}
	}

	// add/set the key & value to the entity form in order
	for f, ov := range fieldOrderValue {
		for i := 0; i < len(ov); i++ {
			position := fmt.Sprintf("%d", i)
			fieldValue := ov[position]

			if payload.Get(f) == "" {
				for i, fv := range fieldValue {
					if i == 0 {
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

}
