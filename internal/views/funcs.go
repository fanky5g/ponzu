package views

import (
	"errors"
	"html/template"
	"time"
)

var GlobFuncs = template.FuncMap{
	"dict": func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("invalid call")
		}

		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, errors.New("dict keys must be strings")
			}

			dict[key] = values[i+1]
		}

		return dict, nil
	},
	"subtract": func(a, b int) int {
		return a - b
	},
	"multiply": func(a, b int) int {
		return a * b
	},
	"formatTime": func(t int64) string {
		return time.Unix(t/1000, 0).Format("02.01.06 03:04 PM")
	},
}
