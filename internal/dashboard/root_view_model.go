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

func NewDashboardRootViewModel(propCache config.ApplicationPropertiesCache) (*DashboardRootViewModel, error) {
	appName, err := propCache.GetAppName()
	if err != nil {
		return nil, err
	}

	publicPath, err := propCache.GetPublicPath()
	if err != nil {
		return nil, err
	}

	contentTypes, err := propCache.GetContentTypes()
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
