package workflow

import "errors"

var ErrWorkflowIncorrectlyConfigured = errors.New("Workflow is incorrectly configured")

type Workflow interface {
	GetState() State
	GetValidTransitions() []Workflow
	GetPastTransitions() []Workflow
}

func GetRoot(workflows []Workflow) (Workflow, error) {
	if len(workflows) == 0 {
		return nil, ErrWorkflowIncorrectlyConfigured
	}

	var rootWorkflow Workflow
	for _, workflow := range workflows {
		if len(workflow.GetPastTransitions()) == 0 {
			if rootWorkflow != nil {
				return nil, ErrWorkflowIncorrectlyConfigured
			}

			rootWorkflow = workflow
		}
	}

	return rootWorkflow, nil
}

func GetSystemDefault() []Workflow {
	return []Workflow{
		new(Draft),
		new(Preview),
		new(Online),
		new(Offline),
		new(Archived),
	}
}
