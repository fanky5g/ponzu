package content

import (
	"errors"

	"github.com/fanky5g/ponzu/content/workflow"
)

var (
	ErrWorkflowIncorrectlyConfigured = errors.New("Workflow is incorrectly configured")
	ErrInvalidWorkflowState          = errors.New("Invalid workflow state")
	ErrWorkflowUnsupported           = errors.New("Workflow not supported for selected content type")
	ErrWorkflowTransitionFailed      = errors.New("Workflow transition failed")
)

func (s *Service) TransitionWorkflowState(entityType, entityId string, targetState workflow.State) (interface{}, error) {
	if _, err := s.Type(entityType); err != nil {
		return nil, err
	}

	entity, err := s.GetContent(entityType, entityId)
	if err != nil {
		return nil, err
	}

	workflowLifecycleSupportedEntity, ok := entity.(workflow.LifecycleSupportedEntity)
	if !ok {
		return nil, ErrWorkflowUnsupported
	}

	if err := s.transitionWorkflowState(workflowLifecycleSupportedEntity, targetState); err != nil {
		return nil, err
	}

	return s.UpdateContent(entityType, entityId, entity)
}

func (s *Service) transitionWorkflowState(entity workflow.LifecycleSupportedEntity, targetState workflow.State) error {
	target := targetState.ToWorkflow()
	if target == nil {
		return ErrInvalidWorkflowState
	}

	w, err := getContentWorkflow(entity)
	if err != nil {
		return err
	}

	for _, tw := range w.GetValidTransitions() {
		if tw == target {
			if err = setContentWorkflow(entity, target); err != nil {
				return err
			}

			break
		}
	}

	cw, err := getContentWorkflow(entity)
	if err != nil {
		return err
	}

	if cw == nil {
		return ErrWorkflowUnsupported
	}

	if cw.GetState() != targetState {
		return ErrWorkflowTransitionFailed
	}

	return nil
}

func getContentWorkflow(entity interface{}) (workflow.Workflow, error) {
	workflowSupportedEntity, ok := entity.(workflow.LifecycleSupportedEntity)
	if !ok {
		return nil, ErrWorkflowUnsupported
	}

	var currentWorkflow workflow.Workflow
	currentWorkflowState := workflowSupportedEntity.GetState()
	if currentWorkflowState != "" {
		currentWorkflow = currentWorkflowState.ToWorkflow()
	}

	if currentWorkflow == nil {
		supportedWorkflows := workflowSupportedEntity.GetSupportedWorkflows()

		var err error
		currentWorkflow, err = getRootWorkflow(supportedWorkflows)
		if err != nil {
			return nil, err
		}
	}

	return currentWorkflow, nil
}

func setContentWorkflow(entity interface{}, w workflow.Workflow) error {
	workflowStateManager, ok := entity.(workflow.LifecycleSupportedEntity)
	if !ok {
		return ErrWorkflowUnsupported
	}

	workflowStateManager.SetState(w.GetState())
	return nil
}

func getRootWorkflow(workflows []workflow.Workflow) (workflow.Workflow, error) {
	if len(workflows) == 0 {
		return nil, ErrWorkflowIncorrectlyConfigured
	}

	var rootWorkflow workflow.Workflow
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
