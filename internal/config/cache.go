package config

import (
	"errors"
	"reflect"

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

type Cache struct {
	cache cache.Cache
}

const CacheKey = "config"

func (c *Cache) GetAppName() (string, error) {
	return c.getConfigString(AppNameJSONKey)
}

func (c *Cache) GetHTTPCacheDisabled() (bool, error) {
	return c.getConfigBool(CacheDisabledJSONKey)
}

func (c *Cache) GetCacheControlMaxAge() (int64, error) {
	return c.getConfigInt(CacheMaxAgeJSONKey)
}

func (c *Cache) GetETag() (string, error) {
	return c.getConfigString(ETagJSONKey)
}

func (c *Cache) GetGZipDisabled() (bool, error) {
	return c.getConfigBool(GzipDisabledJSONKey)
}

func (c *Cache) GetCorsDisabled() (bool, error) {
	return c.getConfigBool(CorsDisabledJSONKey)
}

func (c *Cache) GetDomain() (string, error) {
	return c.getConfigString(DomainJSONKey)
}

func (c *Cache) getConfigString(key string) (string, error) {
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

func (c *Cache) getConfigInt(key string) (int64, error) {
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

func (c *Cache) getConfigBool(key string) (bool, error) {
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

func (c *Cache) getConfigValue(key string) (*reflect.Value, error) {
	cacheValue := c.cache.Get(CacheKey)
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

func NewCache(cache cache.Cache) (*Cache, error) {
	return &Cache{
		cache: cache,
	}, nil
}
