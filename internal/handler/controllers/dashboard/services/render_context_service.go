package services

import (
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/handler/controllers/dashboard/resources"
)

func GetRootRenderContext(propCache config.ApplicationPropertiesCache) (*resources.RootRenderContext, error) {
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

	return &resources.RootRenderContext{
		PublicPath: publicPath,
		AppName:    appName,
		Logo:       appName,
		Types:      contentTypes,
	}, nil
}
