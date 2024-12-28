package content

import (
	"testing"

	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WorkflowTestSuite struct {
	suite.Suite
}

type story struct {
	item.Item
}

func (s *story) GetSupportedWorkflows() []workflow.Workflow {
	return workflow.GetSystemDefault()
}

func (suite *WorkflowTestSuite) TestTransitionWorkflow() {
	folklore := &story{
		item.Item{
			WorkflowState: workflow.DraftState,
		},
	}

	service := &Service{}
	assert.NoError(suite.T(), service.transitionWorkflowState(folklore, workflow.PreviewState))
	assert.Equal(suite.T(), workflow.PreviewState, folklore.GetState())
}

func (suite *WorkflowTestSuite) TestTransitionWorkflowFails() {
	folklore := &story{
		item.Item{
			WorkflowState: workflow.DraftState,
		},
	}

	service := &Service{}
	err := service.transitionWorkflowState(folklore, workflow.OnlineState)
	assert.EqualError(suite.T(), err, ErrWorkflowTransitionFailed.Error())
}

func TestWorkflow(t *testing.T) {
	suite.Run(t, new(WorkflowTestSuite))
}
