package workflow

import (
	"bytes"
	"html/template"
	"strings"
)

var TemplateFuncs = template.FuncMap{
	"WorkflowActionName": func(transition Workflow, currentWorkflow Workflow, entity interface{}) (string, error) {
		action := string(transition.GetState())

		if actionDescriptor, ok := transition.(ActionDescriptor); ok {
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
	"WorkflowStateToLower": func(state State) string {
		return strings.ToLower(string(state))
	},
}
