package content

import (
	"errors"

	"github.com/fanky5g/ponzu/content/workflow"
)

var (
	ErrWorkflowIncorrectlyConfigured = errors.New("workflow is incorrectly configured")
	ErrInvalidWorkflowState          = errors.New("invalid workflow state")
	ErrWorkflowUnsupported           = errors.New("workflow not supported for selected content type")
	ErrWorkflowTransitionFailed      = errors.New("workflow transition failed")
)

func (s *Service) TransitionWorkflowState(entityType, entityId string, targetState workflow.State) (interface{}, error) {
	if _, typeErr := s.Type(entityType); typeErr != nil {
		return nil, typeErr
	}

	entity, getContentErr := s.GetContent(entityType, entityId)
	if getContentErr != nil {
		return nil, getContentErr
	}

	workflowEntity, ok := entity.(workflow.LifecycleSupportedEntity)
	if !ok {
		return nil, ErrWorkflowUnsupported
	}

	currentState := workflowEntity.GetState()
	if transitionWorkflowStateErr := s.transitionWorkflowState(workflowEntity, targetState); transitionWorkflowStateErr != nil {
		return nil, transitionWorkflowStateErr
	}

	updated, updateContentErr := s.UpdateContent(entityType, entityId, entity)
	if updateContentErr != nil {
		return nil, updateContentErr
	}

	if workflowStateChangeErr := s.notifyWorkflowStateChangeListeners(
		currentState,
		updated.(workflow.LifecycleSupportedEntity),
	); workflowStateChangeErr != nil {
		// TODO(B.B): remove after introducing (unit of work concept - transactions).
		updatedWorkflow, getContentWorkflowErr := getContentWorkflow(updated)
		if getContentWorkflowErr != nil {
			return nil, getContentWorkflowErr
		}

		if updatedWorkflow.GetState() == targetState {
			err := setContentWorkflow(updated, currentState.ToWorkflow())
			if err != nil {
				return nil, err
			}

			reverted, revertErr := s.UpdateContent(entityType, entityId, updated)
			if revertErr != nil {
				return nil, revertErr
			}

			return reverted, workflowStateChangeErr
		}

		return updated, workflowStateChangeErr
	}

	return updated, updateContentErr
}

func (s *Service) notifyWorkflowStateChangeListeners(prevState workflow.State, entity workflow.LifecycleSupportedEntity) error {
	referenceLoader := newReferenceLoader(entity, s)
	if entityWorkflowStateChangeTrigger, ok := entity.(workflow.EntityStateChangeTrigger); ok {
		if workflowStateChangeTriggerErr := entityWorkflowStateChangeTrigger.OnWorkflowStateChange(
			prevState,
			referenceLoader,
		); workflowStateChangeTriggerErr != nil {
			return workflowStateChangeTriggerErr
		}
	}

	if s.workflowStateChangeHandler != nil {
		return s.workflowStateChangeHandler.OnWorkflowStateChange(
			entity,
			prevState,
			referenceLoader,
		)
	}

	return nil
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
	for _, w := range workflows {
		if len(w.GetPastTransitions()) == 0 {
			if rootWorkflow != nil {
				return nil, ErrWorkflowIncorrectlyConfigured
			}

			rootWorkflow = w
		}
	}

	return rootWorkflow, nil
}
