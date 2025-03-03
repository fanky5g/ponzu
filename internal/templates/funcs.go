package templates

import (
	"bytes"
	"errors"
	"github.com/fanky5g/ponzu/content/workflow"
	"html/template"
	"strings"

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
	"FmtTime":  util.FmtTime,
	"FmtBytes": util.FmtBytes,
	"WorkflowActionName": func(transition workflow.Workflow, currentWorkflow workflow.Workflow, entity interface{}) (string, error) {
		action := string(transition.GetState())

		if actionDescriptor, ok := transition.(workflow.ActionDescriptor); ok {
			actionTemplate, err := actionDescriptor.GetAction(currentWorkflow)
			if err != nil {
				return "", err
			}

			w := &bytes.Buffer{}
			if err = actionTemplate.Execute(w, entity); err != nil {
				return "", err
			}

			action = w.String()
		}

		return action, nil
	},
	"WorkflowStateToLower": func(state workflow.State) string {
		return strings.ToLower(string(state))
	},
}
