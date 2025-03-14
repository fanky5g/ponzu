package dashboard

import (
	"github.com/fanky5g/ponzu/content"
)

type ViewModel struct {
	PublicPath string
	AppName    string
	Logo       string
	Types      map[string]content.Builder
}

func NewDashboardViewModel(
	appNameProvider AppNameProvider,
	publicPath string,
	contentTypes map[string]content.Builder,
) (*ViewModel, error) {
	appName, err := appNameProvider.GetAppName()
	if err != nil {
		return nil, err
	}

	return &ViewModel{
		PublicPath: publicPath,
		AppName:    appName,
		Logo:       appName,
		Types:      contentTypes,
	}, nil
}
