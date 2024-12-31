package uploads

import (
	"html/template"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/dashboard"
)

type EditUploadFormViewModel struct {
	dashboard.DashboardRootViewModel
	ID         string
	EntityName string
	Slug       string
	Form       template.HTML
}

// TODO: upload is a valid content type. refactor common content mapping
func NewEditUploadFormViewModel(
	upload *entities.FileUpload,
	cfg config.ConfigCache,
	publicPath string,
	contentTypes map[string]content.Builder) (*EditUploadFormViewModel, error) {
	formBytes, err := upload.MarshalEditor(publicPath)
	if err != nil {
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

	return &EditUploadFormViewModel{
		DashboardRootViewModel: *rootViewModel,
		ID:                     upload.ID,
		EntityName:             upload.EntityName(),
		Slug:                   upload.Slug,
		Form:                   template.HTML(formBytes),
	}, nil
}
