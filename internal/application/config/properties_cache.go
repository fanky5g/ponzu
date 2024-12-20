package config

import (
	"errors"

	"github.com/fanky5g/ponzu/cache"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/entities"
	"github.com/fanky5g/ponzu/util"
)

type propCache struct {
	cache        cache.Cache
	cfg          *config.Config
	contentTypes map[string]content.Builder
}

const ConfigCacheKey = "config"

func (c *propCache) GetAppName() (string, error) {
	cacheValue := c.cache.Get(ConfigCacheKey)
	configCache, ok := cacheValue.(*entities.Config)
	if !ok {
		return "", errors.New("missing application config")
	}

	return util.StringFieldByJSONTagName(configCache, "name")
}

func (c *propCache) GetPublicPath() (string, error) {
	return c.cfg.Paths.PublicPath, nil
}

func (c *propCache) GetContentTypes() (map[string]content.Builder, error) {
	return c.contentTypes, nil
}

func NewApplicationPropertiesCache(
	cache cache.Cache,
	contentTypes map[string]content.Builder,
) (*propCache, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	return &propCache{
		cache:        cache,
		cfg:          cfg,
		contentTypes: contentTypes,
	}, nil
}
