package mocks

import (
	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/stretchr/testify/mock"
)

type WorkflowStateChangeHandler struct {
	*mock.Mock
}

func (handler *WorkflowStateChangeHandler) OnWorkflowStateChange(
	entity workflow.LifecycleSupportedEntity,
	prevState workflow.State,
	referenceLoader workflow.EntityLoader,
) error {
	args := handler.Mock.Called(entity, prevState, referenceLoader)
	return args.Error(0)
}
