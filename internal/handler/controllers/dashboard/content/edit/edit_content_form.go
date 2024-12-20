package edit

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/fanky5g/ponzu/internal/handler/controllers/dashboard/resources"
	"github.com/fanky5g/ponzu/internal/handler/controllers/dashboard/services"
)

var ErrInvalidContentType = errors.New("invalid content type")

type ContentMetadata struct {
	ID         string
	EntityName string
	Slug       string
	Title      string
	Workflow   workflow.Workflow
}

type EditContentForm struct {
	resources.RootRenderContext
	ContentMetadata
	Form template.HTML
}

func NewEditContentForm(
	entity interface{},
	propCache config.ApplicationPropertiesCache,
) (*EditContentForm, error) {
	entityInterface, ok := entity.(content.Entity)
	if !ok {
		return nil, ErrInvalidContentType
	}

	editable, ok := entity.(editor.Editable)
	if !ok {
		return nil, fmt.Errorf("entities type %T is not editable", entityInterface.EntityName())
	}

	rootRenderContext, err := services.GetRootRenderContext(propCache)
	if err != nil {
		return nil, err
	}

	formBytes, err := editable.MarshalEditor(rootRenderContext.PublicPath)
	if err != nil {
		return nil, err
	}

	identifiable, ok := entity.(item.Identifiable)
	if !ok {
		return nil, ErrInvalidContentType
	}

	slug := ""
	if sluggable, ok := entity.(item.Sluggable); ok {
		slug = sluggable.ItemSlug()
	}

	itemTitle := entityInterface.EntityName()
	if identifiable, ok := entity.(item.Readable); ok {
		itemTitle = identifiable.GetTitle()
	}

	currentWorkflow, err := workflow.GetContentWorkflow(entity)
	if err != nil {
		return nil, err
	}

	return &EditContentForm{
		RootRenderContext: *rootRenderContext,
		ContentMetadata: ContentMetadata{
			ID:         identifiable.ItemID(),
			EntityName: entityInterface.EntityName(),
			Title:      itemTitle,
			Slug:       slug,
			Workflow:   currentWorkflow,
		},
		Form: template.HTML(formBytes),
	}, nil
}
