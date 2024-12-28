package workflow

import (
	"text/template"
)

type Workflow interface {
	GetState() State
	GetValidTransitions() []Workflow
	GetPastTransitions() []Workflow
}

type LifecycleSupportedEntity interface {
	GetSupportedWorkflows() []Workflow
	SetState(state State)
	GetState() State
}

type ActionDescriptor interface {
	GetAction(source Workflow) (*template.Template, error)
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
