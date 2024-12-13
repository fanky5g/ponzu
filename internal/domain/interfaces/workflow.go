package interfaces

import "github.com/fanky5g/ponzu/content/workflow"

type MutableWorkflowState interface {
	SetState(state workflow.State)
}
