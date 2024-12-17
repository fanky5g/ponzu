package edit

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/fanky5g/ponzu/internal/domain/interfaces"
	"github.com/fanky5g/ponzu/internal/handler/controllers/dashboard/resources"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	configService "github.com/fanky5g/ponzu/internal/services/config"
	contentService "github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/tokens"
)

type WorkflowAction struct {
	State  workflow.State
	Action string
}

type WorkflowMetadata struct {
	State        workflow.State
	TargetStates []WorkflowAction
}

type ContentMetadata struct {
	ID       string
	Kind     string
	Slug     string
	Title    string
	Workflow *WorkflowMetadata
}

type RootRenderContext struct {
	PublicPath string
	AppName    string
	Logo       string
	Types      map[string]content.Builder
}

type RenderContext struct {
	RootRenderContext
	ContentMetadata
	Form template.HTML
}

type RenderContextService interface {
	GetEditorRenderContext(identifier *resources.ContentIdentifier) (*RenderContext, error)
}

type renderContextService struct {
	contentService contentService.Service
	configService  configService.Service
	paths          config.Paths
	Types          content.Types
}

func (s *renderContextService) GetEditorRenderContext(identifier *resources.ContentIdentifier) (*RenderContext, error) {
	typeBuilder, ok := s.Types.Content[identifier.Type]
	if !ok {
		return nil, content.ErrTypeNotRegistered
	}

	entity := typeBuilder()
	var err error
	if identifier.ID != "" {
		entity, err = s.contentService.GetContent(identifier.Type, identifier.ID)
		if err != nil {
			return nil, err
		}
	}

	editable, ok := entity.(editor.Editable)
	if !ok {
		return nil, fmt.Errorf("entities type %T is not editable", identifier.Type)
	}

	formBytes, err := editable.MarshalEditor(s.paths)
	if err != nil {
		return nil, err
	}

	slug := ""
	if sluggable, ok := entity.(item.Sluggable); ok {
		slug = sluggable.ItemSlug()
	}

	itemTitle := identifier.Type
	if identifiable, ok := entity.(item.Readable); ok {
		itemTitle = identifiable.GetTitle()
	}

	rootRenderContext, err := s.GetRootRenderContext()
	if err != nil {
		return nil, err
	}

	workflowMetadata, err := s.GetWorkflowMetadata(entity)
	if err != nil {
		return nil, err
	}

	return &RenderContext{
		RootRenderContext: *rootRenderContext,
		ContentMetadata: ContentMetadata{
			ID:       identifier.ID,
			Kind:     identifier.Type,
			Title:    itemTitle,
			Slug:     slug,
			Workflow: workflowMetadata,
		},
		Form: template.HTML(formBytes),
	}, nil
}

func (s *renderContextService) GetRootRenderContext() (*RootRenderContext, error) {
	appName, err := s.configService.GetAppName()
	if err != nil {
		return nil, err
	}

	return &RootRenderContext{
		PublicPath: s.paths.PublicPath,
		AppName:    appName,
		Logo:       appName,
		Types:      s.Types.Content,
	}, nil
}

func (s *renderContextService) GetWorkflowMetadata(entity interface{}) (*WorkflowMetadata, error) {
	workflowLifecycle, ok := entity.(interfaces.HasWorkflowLifecycle)
	if !ok {
		return nil, nil
	}

	var currentWorkflow workflow.Workflow
	if workflowStateManager, ok := entity.(interfaces.WorkflowStateManager); ok {
		currentWorkflowState := workflowStateManager.GetState()
		if currentWorkflowState != "" {
			currentWorkflow = currentWorkflowState.ToWorkflow()
		}
	}

	if currentWorkflow == nil {
		supportedWorkflows := workflowLifecycle.GetSupportedWorkflows()

		var err error
		currentWorkflow, err = workflow.GetRoot(supportedWorkflows)
		if err != nil {
			return nil, err
		}
	}

	validTransitions := currentWorkflow.GetValidTransitions()
	targetWorkflowStates := make([]WorkflowAction, len(validTransitions))
	for i, transition := range validTransitions {
		state := transition.GetState()
		action := string(state)

		if actionDescriptor, ok := transition.(interfaces.WorkflowActionDescriptor); ok {
			actionTemplate, err := actionDescriptor.GetAction(currentWorkflow)
			if err != nil {
				return nil, err
			}

			w := &bytes.Buffer{}
			if err = actionTemplate.Execute(w, entity); err != nil {
				return nil, err
			}

			action = w.String()
		}

		targetWorkflowStates[i] = WorkflowAction{
			State:  state,
			Action: action,
		}
	}

	return &WorkflowMetadata{
		State:        currentWorkflow.GetState(),
		TargetStates: targetWorkflowStates,
	}, nil
}

func NewRenderContextService(
	r router.Router,
) (RenderContextService, error) {
	configService := r.Context().Service(tokens.ConfigServiceToken).(configService.Service)
	contentService := r.Context().Service(tokens.ContentServiceToken).(contentService.Service)

	return &renderContextService{
		contentService: contentService,
		configService:  configService,
		paths:          r.Context().Paths(),
		Types:          r.Context().Types(),
	}, nil
}
