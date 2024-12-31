package dashboard

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/config"
)

type DashboardRootViewModel struct {
	PublicPath string
	AppName    string
	Logo       string
	Types      map[string]content.Builder
}

func NewDashboardRootViewModel(
	cfg config.ConfigCache,
	publicPath string,
	contentTypes map[string]content.Builder,
) (*DashboardRootViewModel, error) {
	appName, err := cfg.GetAppName()
	if err != nil {
		return nil, err
	}

	return &DashboardRootViewModel{
		PublicPath: publicPath,
		AppName:    appName,
		Logo:       appName,
		Types:      contentTypes,
	}, nil
}
