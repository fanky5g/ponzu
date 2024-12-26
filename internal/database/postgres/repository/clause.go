package repository

import (
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/internal/constants"
	"reflect"
	"time"
)

func (repo *Repository) valueType(value interface{}) (string, error) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Interface().(type) {
	case string:
		return "text", nil
	case bool:
		return "boolean", nil
	case int:
		return "numeric", nil
	case int64:
		return "numeric", nil
	case int32:
		return "numeric", nil
	case int8:
		return "numeric", nil
	case int16:
		return "numeric", nil
	case float64:
		return "numeric", nil
	case float32:
		return "numeric", nil
	case time.Time:
		return "timestamptz", nil
	default:
		return "", fmt.Errorf("unsupported type %T", value)
	}
}

func (repo *Repository) getComparisonOperator(operator constants.ComparisonOperator) (string, error) {
	var comparisonOperator string

	switch operator {
	case constants.Equal:
		comparisonOperator = "="
	case constants.LessThan:
		comparisonOperator = "<"
	case constants.LessThanOrEqualTo:
		comparisonOperator = "<="
	case constants.GreaterThan:
		comparisonOperator = ">"
	case constants.GreaterThanOrEqualTo:
		comparisonOperator = ">="
	default:
		return "", errors.New("invalid comparison operator")
	}

	return comparisonOperator, nil
}
