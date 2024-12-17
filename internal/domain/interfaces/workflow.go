package interfaces

import (
	"text/template"

	"github.com/fanky5g/ponzu/content/workflow"
)

type HasWorkflowLifecycle interface {
	GetSupportedWorkflows() []workflow.Workflow
}

type WorkflowStateManager interface {
	SetState(state workflow.State)
	GetState() workflow.State
}

type WorkflowActionDescriptor interface {
	GetAction(source workflow.Workflow) (*template.Template, error)
}
