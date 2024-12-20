package config

import "github.com/fanky5g/ponzu/content"

type ApplicationPropertiesCache interface {
	GetAppName() (string, error)
	GetPublicPath() (string, error)
	GetContentTypes() (map[string]content.Builder, error)
}
