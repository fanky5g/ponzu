package config

import (
	"errors"
	"reflect"

	globConfig "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/cache"
	"github.com/fanky5g/ponzu/util"
	log "github.com/sirupsen/logrus"
)

var (
	ErrInvalidConfig            = errors.New("invalid configuration")
	ErrApplicationConfigMissing = errors.New("missing application config")
	CacheDisabledJSONKey        = "cache_disabled"
	AppNameJSONKey              = "name"
	CacheMaxAgeJSONKey          = "cache_max_age"
	ETagJSONKey                 = "etag"
	CorsDisabledJSONKey         = "cors_disabled"
	GzipDisabledJSONKey         = "gzip_disabled"
	DomainJSONKey               = "domain"
)

type propCache struct {
	cache        cache.Cache
	cfg          *globConfig.Config
	contentTypes map[string]content.Builder
}

type ApplicationPropertiesCache interface {
	GetAppName() (string, error)
	GetHTTPCacheDisabled() (bool, error)
	GetETag() (string, error)
	GetCacheControlMaxAge() (int64, error)
	GetGZipDisabled() (bool, error)
	GetCorsDisabled() (bool, error)
	GetDomain() (string, error)
	GetPublicPath() (string, error)
	GetContentTypes() (map[string]content.Builder, error)
}

const ConfigCacheKey = "config"

func (c *propCache) GetAppName() (string, error) {
	return c.getConfigString(AppNameJSONKey)
}

func (c *propCache) GetHTTPCacheDisabled() (bool, error) {
	return c.getConfigBool(CacheDisabledJSONKey)
}

func (c *propCache) GetCacheControlMaxAge() (int64, error) {
	return c.getConfigInt(CacheMaxAgeJSONKey)
}

func (c *propCache) GetETag() (string, error) {
	return c.getConfigString(ETagJSONKey)
}

func (c *propCache) GetGZipDisabled() (bool, error) {
	return c.getConfigBool(GzipDisabledJSONKey)
}

func (c *propCache) GetCorsDisabled() (bool, error) {
	return c.getConfigBool(CorsDisabledJSONKey)
}

func (c *propCache) GetDomain() (string, error) {
	return c.getConfigString(DomainJSONKey)
}

func (c *propCache) getConfigString(key string) (string, error) {
	value, err := c.getConfigValue(key)
	if err != nil {
		return "", err
	}

	if value.Kind() != reflect.String {
		log.Warnf("%s is not a valid string value", key)
		return "", ErrInvalidConfig
	}

	return value.String(), nil
}

func (c *propCache) getConfigInt(key string) (int64, error) {
	value, err := c.getConfigValue(key)
	if err != nil {
		return 0, err
	}

	if value.Kind() != reflect.Int64 {
		log.Warnf("%s is not a valid int64 value", key)
		return 0, ErrInvalidConfig
	}

	return value.Int(), nil
}

func (c *propCache) getConfigBool(key string) (bool, error) {
	value, err := c.getConfigValue(key)
	if err != nil {
		return false, err
	}

	if value.Kind() != reflect.Bool {
		log.Warnf("%s is not a valid boolean value", key)
		return false, ErrInvalidConfig
	}

	return value.Bool(), nil
}

func (c *propCache) getConfigValue(key string) (*reflect.Value, error) {
	cacheValue := c.cache.Get(ConfigCacheKey)
	cfg, ok := cacheValue.(*entities.Config)
	if !ok {
		return nil, ErrApplicationConfigMissing
	}

	value := util.FieldByJSONTagName(cfg, key)
	if !value.IsValid() {
		log.Warnf("%s is not a valid config entry", key)
		return nil, ErrInvalidConfig
	}

	return &value, nil
}

func (c *propCache) GetPublicPath() (string, error) {
	return c.cfg.Paths.PublicPath, nil
}

func (c *propCache) GetContentTypes() (map[string]content.Builder, error) {
	return c.contentTypes, nil
}

func NewApplicationPropertiesCache(cache cache.Cache, contentTypes map[string]content.Builder) (*propCache, error) {
	cfg, err := globConfig.Get()
	if err != nil {
		return nil, err
	}

	return &propCache{
		cache:        cache,
		cfg:          cfg,
		contentTypes: contentTypes,
	}, nil
}
