package config

import (
	"errors"
	"reflect"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/cache"
	"github.com/fanky5g/ponzu/util"
	log "github.com/sirupsen/logrus"
)

var (
	ErrInvalidConfig     = errors.New("invalid configuration")
	CacheDisabledJSONKey = "cache_disabled"
	AppNameJSONKey       = "name"
	CacheMaxAgeJSONKey   = "cache_max_age"
	ETagJSONKey          = "etag"
	CorsDisabledJSONKey  = "cors_disabled"
	GzipDisabledJSONKey  = "gzip_disabled"
	DomainJSONKey        = "domain"
)

type configCache struct {
	cache cache.Cache
}

type ConfigCache interface {
	GetAppName() (string, error)
	GetHTTPCacheDisabled() (bool, error)
	GetETag() (string, error)
	GetCacheControlMaxAge() (int64, error)
	GetGZipDisabled() (bool, error)
	GetCorsDisabled() (bool, error)
	GetDomain() (string, error)
}

const ConfigCacheKey = "config"

func (c *configCache) GetAppName() (string, error) {
	return c.getConfigString(AppNameJSONKey)
}

func (c *configCache) GetHTTPCacheDisabled() (bool, error) {
	return c.getConfigBool(CacheDisabledJSONKey)
}

func (c *configCache) GetCacheControlMaxAge() (int64, error) {
	return c.getConfigInt(CacheMaxAgeJSONKey)
}

func (c *configCache) GetETag() (string, error) {
	return c.getConfigString(ETagJSONKey)
}

func (c *configCache) GetGZipDisabled() (bool, error) {
	return c.getConfigBool(GzipDisabledJSONKey)
}

func (c *configCache) GetCorsDisabled() (bool, error) {
	return c.getConfigBool(CorsDisabledJSONKey)
}

func (c *configCache) GetDomain() (string, error) {
	return c.getConfigString(DomainJSONKey)
}

func (c *configCache) getConfigString(key string) (string, error) {
	value, err := c.getConfigValue(key)
	if err != nil {
		return "", err
	}

	if value == nil {
		return "", nil
	}

	if value.Kind() != reflect.String {
		log.Warnf("%s is not a valid string value", key)
		return "", ErrInvalidConfig
	}

	return value.String(), nil
}

func (c *configCache) getConfigInt(key string) (int64, error) {
	value, err := c.getConfigValue(key)
	if err != nil {
		return 0, err
	}

	if value == nil {
		return 0, nil
	}

	if value.Kind() != reflect.Int64 {
		log.Warnf("%s is not a valid int64 value", key)
		return 0, ErrInvalidConfig
	}

	return value.Int(), nil
}

func (c *configCache) getConfigBool(key string) (bool, error) {
	value, err := c.getConfigValue(key)
	if err != nil {
		return false, err
	}

	if value == nil {
		return false, nil
	}

	if value.Kind() != reflect.Bool {
		log.Warnf("%s is not a valid boolean value", key)
		return false, ErrInvalidConfig
	}

	return value.Bool(), nil
}

func (c *configCache) getConfigValue(key string) (*reflect.Value, error) {
	cacheValue := c.cache.Get(ConfigCacheKey)
	cfg, ok := cacheValue.(*Config)
	if !ok {
		return nil, nil
	}

	value := util.FieldByJSONTagName(cfg, key)
	if !value.IsValid() {
		log.Warnf("%s is not a valid config entry", key)
		return nil, ErrInvalidConfig
	}

	return &value, nil
}

func NewCache(cache cache.Cache, contentTypes map[string]content.Builder) (*configCache, error) {
	return &configCache{
		cache: cache,
	}, nil
}
