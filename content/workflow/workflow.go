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

type EntityLoader interface {
	GetEntity(entityName, entityId string) (interface{}, error)
}

type EntityStateChangeTrigger interface {
	OnWorkflowStateChange(prevState State, entityLoader EntityLoader) error
}

type StateChangeTrigger interface {
	OnWorkflowStateChange(entity LifecycleSupportedEntity, prevState State, entityLoader EntityLoader) error
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
