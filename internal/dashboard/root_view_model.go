package dashboard

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/config"
)

type RootViewModel struct {
	PublicPath string
	AppName    string
	Logo       string
	Types      map[string]content.Builder
}

func NewDashboardRootViewModel(
	cfg config.ConfigCache,
	publicPath string,
	contentTypes map[string]content.Builder,
) (*RootViewModel, error) {
	appName, err := cfg.GetAppName()
	if err != nil {
		return nil, err
	}

	return &RootViewModel{
		PublicPath: publicPath,
		AppName:    appName,
		Logo:       appName,
		Types:      contentTypes,
	}, nil
}
