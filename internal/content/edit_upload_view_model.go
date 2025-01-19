package content

import (
	"html/template"

	"github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/internal/dashboard"
)

type EditUploadFormViewModel struct {
	dashboard.RootViewModel
	ID         string
	EntityName string
	Slug       string
	Form       template.HTML
}

func NewEditUploadFormViewModel(
	upload *entities.Upload,
	rootViewModel *dashboard.RootViewModel,
) (*EditUploadFormViewModel, error) {
	formBytes, err := upload.MarshalEditor(rootViewModel.PublicPath)
	if err != nil {
		return nil, err
	}

	return &EditUploadFormViewModel{
		RootViewModel: *rootViewModel,
		ID:            upload.ID,
		EntityName:    upload.EntityName(),
		Slug:          upload.Slug,
		Form:          template.HTML(formBytes),
	}, nil
}
