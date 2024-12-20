package content

import (
	"github.com/fanky5g/ponzu/internal/config"
)

func GetRootRenderContext(propCache config.ApplicationPropertiesCache) (*RootRenderContext, error) {
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

	return &RootRenderContext{
		PublicPath: publicPath,
		AppName:    appName,
		Logo:       appName,
		Types:      contentTypes,
	}, nil
}
