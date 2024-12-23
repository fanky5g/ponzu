package workflow

import (
	"errors"
	"text/template"
)

var ErrWorkflowIncorrectlyConfigured = errors.New("Workflow is incorrectly configured")

type Workflow interface {
	GetState() State
	GetValidTransitions() []Workflow
	GetPastTransitions() []Workflow
}

type HasWorkflowLifecycle interface {
	GetSupportedWorkflows() []Workflow
}

type StateManager interface {
	SetState(state State)
	GetState() State
}

type ActionDescriptor interface {
	GetAction(source Workflow) (*template.Template, error)
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

func GetContentWorkflow(entity interface{}) (Workflow, error) {
	workflowLifecycle, ok := entity.(HasWorkflowLifecycle)
	if !ok {
		return nil, nil
	}

	var currentWorkflow Workflow
	if workflowStateManager, ok := entity.(StateManager); ok {
		currentWorkflowState := workflowStateManager.GetState()
		if currentWorkflowState != "" {
			currentWorkflow = currentWorkflowState.ToWorkflow()
		}
	}

	if currentWorkflow == nil {
		supportedWorkflows := workflowLifecycle.GetSupportedWorkflows()

		var err error
		currentWorkflow, err = GetRoot(supportedWorkflows)
		if err != nil {
			return nil, err
		}
	}

	return currentWorkflow, nil
}
