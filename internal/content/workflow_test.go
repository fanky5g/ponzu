package content

import (
	"errors"
	"testing"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/fanky5g/ponzu/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WorkflowTestSuite struct {
	suite.Suite
	service *Service
	m       *mock.Mock
}

type story struct {
	item.Item
}

func (s *story) GetSupportedWorkflows() []workflow.Workflow {
	return workflow.GetSystemDefault()
}

func (s *story) GetRepositoryToken() string {
	return "Story"
}

type storyWithWorkflowStateChangeTrigger struct {
	item.Item
	m *mock.Mock
}

func (s *storyWithWorkflowStateChangeTrigger) GetSupportedWorkflows() []workflow.Workflow {
	return workflow.GetSystemDefault()
}

func (s *storyWithWorkflowStateChangeTrigger) GetRepositoryToken() string {
	return "Story"
}

func (s *storyWithWorkflowStateChangeTrigger) OnWorkflowStateChange(prevState workflow.State) error {
	args := s.m.Called(prevState)
	return args.Error(0)
}

func (suite *WorkflowTestSuite) SetupSuite() {
	contentTypes := map[string]content.Builder{
		"Story": func() interface{} {
			return new(story)
		},
	}

	suite.m = &mock.Mock{}
	var err error
	suite.service, err = New(&mocks.DB{Mock: suite.m}, contentTypes, &mocks.SearchClient{Mock: suite.m}, nil)
	if err != nil {
		suite.T().Fatal(err)
		return
	}
}

func (suite *WorkflowTestSuite) TestTransitionWorkflowState() {
	entityType := "Story"
	entityId := "1"

	entity := &story{
		item.Item{
			ID:            entityId,
			WorkflowState: workflow.DraftState,
		},
	}

	update := &story{
		item.Item{
			ID:            entityId,
			WorkflowState: workflow.PreviewState,
		},
	}

	suite.m.On("FindOneById", entityId).Once().Return(entity, nil)
	suite.m.On("UpdateById", entityId, update).Once().Return(update, nil)
	suite.m.On("Update", entityId, update).Once().Return(nil)

	result, err := suite.service.TransitionWorkflowState(entityType, entityId, workflow.PreviewState)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), update, result)
	}
}

func (suite *WorkflowTestSuite) TestTransitionWorkflowStateCallsWorkflowStateChangeTrigger() {
	entityType := "Story"
	entityId := "1"

	entity := &storyWithWorkflowStateChangeTrigger{
		item.Item{
			ID:            entityId,
			WorkflowState: workflow.DraftState,
		},
		suite.m,
	}

	update := &storyWithWorkflowStateChangeTrigger{
		item.Item{
			ID:            entityId,
			WorkflowState: workflow.PreviewState,
		},
		suite.m,
	}

	suite.m.On("FindOneById", entityId).Once().Return(entity, nil)
	suite.m.On("UpdateById", entityId, update).Once().Return(update, nil)
	suite.m.On("Update", entityId, update).Once().Return(nil)
	suite.m.On("OnWorkflowStateChange", workflow.DraftState).Once().Return(nil)

	result, err := suite.service.TransitionWorkflowState(entityType, entityId, workflow.PreviewState)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), update, result)
	}
}

func (suite *WorkflowTestSuite) TestTransitionWorkflowStateReturnsWorkflowStateChangeTriggerError() {
	entityType := "Story"
	entityId := "1"

	entity := &storyWithWorkflowStateChangeTrigger{
		item.Item{
			ID:            entityId,
			WorkflowState: workflow.DraftState,
		},
		suite.m,
	}

	update := &storyWithWorkflowStateChangeTrigger{
		item.Item{
			ID:            entityId,
			WorkflowState: workflow.PreviewState,
		},
		suite.m,
	}

	expectedError := errors.New("Something bad happened")

	suite.m.On("FindOneById", entityId).Once().Return(entity, nil)
	suite.m.On("UpdateById", entityId, update).Return(update, nil).Once()
	suite.m.On("UpdateById", entityId, mock.Anything).Return(entity, nil).Once()
	suite.m.On("Update", entityId, update).Once().Return(nil)
	suite.m.On("Update", entityId, entity).Once().Return(nil)
	suite.m.On("OnWorkflowStateChange", workflow.DraftState).Once().Return(expectedError)

	result, err := suite.service.TransitionWorkflowState(entityType, entityId, workflow.PreviewState)
	assert.EqualError(suite.T(), err, expectedError.Error())
	assert.NotNil(suite.T(), result)
}

func (suite *WorkflowTestSuite) TestTransitionWorkflow() {
	folklore := &story{
		item.Item{
			WorkflowState: workflow.DraftState,
		},
	}

	assert.NoError(suite.T(), suite.service.transitionWorkflowState(folklore, workflow.PreviewState))
	assert.Equal(suite.T(), workflow.PreviewState, folklore.GetState())
}

func (suite *WorkflowTestSuite) TestTransitionWorkflowFails() {
	folklore := &story{
		item.Item{
			WorkflowState: workflow.DraftState,
		},
	}

	err := suite.service.transitionWorkflowState(folklore, workflow.OnlineState)
	assert.EqualError(suite.T(), err, ErrWorkflowTransitionFailed.Error())
}

func TestWorkflow(t *testing.T) {
	suite.Run(t, new(WorkflowTestSuite))
}
