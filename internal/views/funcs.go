package views

import (
	"errors"
	"html/template"
	"strings"

	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/fanky5g/ponzu/util"
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
	"lower": func(a interface{}) string {
		switch a.(type) {
		case string:
			return strings.ToLower(a.(string))
		case workflow.State:
			workflowState, ok := a.(workflow.State)
			if ok {
				return strings.ToLower(string(workflowState))
			}
			return ""
		default:
			return ""
		}
	},
	"FmtTime":  util.FmtTime,
	"FmtBytes": util.FmtBytes,
}
