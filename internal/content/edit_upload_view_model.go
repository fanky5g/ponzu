package content

import (
	"html/template"

	"github.com/fanky5g/ponzu/content/entities"
)

type EditUploadFormViewModel struct {
	PublicPath string
	ID         string
	EntityName string
	Slug       string
	Form       template.HTML
}

func NewEditUploadFormViewModel(publicPath string, upload *entities.Upload) (*EditUploadFormViewModel, error) {
	formBytes, err := upload.MarshalEditor(publicPath)
	if err != nil {
		return nil, err
	}

	return &EditUploadFormViewModel{
		PublicPath: publicPath,
		ID:         upload.ID,
		EntityName: upload.EntityName(),
		Slug:       upload.Slug,
		Form:       template.HTML(formBytes),
	}, nil
}
