package content

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/fanky5g/ponzu/exceptions"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/dashboard"
)

var (
	ErrInvalidContentType = errors.New("invalid content type")
)

type (
	Metadata struct {
		ID         string
		EntityName string
		Slug       string
		Title      string
		Workflow   workflow.Workflow
	}

	EditContentFormViewModel struct {
		dashboard.RootViewModel
		Metadata
		Form  template.HTML
		Error error
	}
)

func NewEditContentFormViewModel(
	entity interface{},
	cfg config.ConfigCache,
	publicPath string,
	contentTypes map[string]content.Builder,
	exception error,
) (*EditContentFormViewModel, error) {
	entityInterface, ok := entity.(content.Entity)
	if !ok {
		return nil, ErrInvalidContentType
	}

	editable, ok := entity.(editor.Editable)
	if !ok {
		return nil, fmt.Errorf("entities type %T is not editable", entityInterface.EntityName())
	}

	formBytes, err := editable.MarshalEditor(publicPath)
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

	currentWorkflow, err := getContentWorkflow(entity)
	if err != nil && !errors.Is(err, ErrWorkflowUnsupported) {
		return nil, err
	}

	rootViewModel, err := dashboard.NewDashboardRootViewModel(
		cfg,
		publicPath,
		contentTypes,
	)

	if err != nil {
		return nil, err
	}

	var clientException *exceptions.ClientException
	if exception != nil {
		errors.As(exception, &clientException)
	}

	return &EditContentFormViewModel{
		RootViewModel: *rootViewModel,
		Metadata: Metadata{
			ID:         identifiable.ItemID(),
			EntityName: entityInterface.EntityName(),
			Title:      itemTitle,
			Slug:       slug,
			Workflow:   currentWorkflow,
		},
		Form:  template.HTML(formBytes),
		Error: clientException,
	}, nil
}
