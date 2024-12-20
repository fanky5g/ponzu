package config

import (
	"fmt"
	"reflect"

	"github.com/fanky5g/ponzu/cache"
	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/entities"
	"github.com/fanky5g/ponzu/util"
)

type Service struct {
	config database.Repository
	cache  cache.Cache
}

func (s *Service) GetStringValue(key string) (string, error) {
	cfg, err := s.config.Latest()
	if err != nil {
		return "", err
	}

	return util.StringFieldByJSONTagName(cfg, key)
}

func (s *Service) GetBoolValue(key string) (bool, error) {
	cfg, err := s.config.Latest()
	if err != nil {
		return false, err
	}

	value := util.FieldByJSONTagName(cfg, key)
	if !value.IsValid() {
		return false, fmt.Errorf("%s is not a valid config entry", key)
	}

	if value.Kind() != reflect.Bool {
		return false, fmt.Errorf("%s is not a boolean", key)
	}

	return value.Bool(), nil
}

func (s *Service) SetConfig(config *entities.Config) error {
	cfg, err := s.Get()
	if err != nil {
		return err
	}

	if cfg == nil {
		_, err = s.config.Insert(config)
		return err
	}

	_, err = s.config.UpdateById(cfg.ID, config)
	if err != nil {
		return err
	}

	s.cache.Set(ConfigCacheKey, config)
	return nil
}

func (s *Service) Get() (*entities.Config, error) {
	cfg, err := s.config.Latest()
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		return nil, nil
	}

	return cfg.(*entities.Config), nil
}

func (s *Service) warmConfigCache() error {
	cfg, err := s.Get()
	if err != nil {
		return err
	}

	if cfg != nil {
		s.cache.Set(ConfigCacheKey, cfg)
	}

	return nil
}

func New(config database.Repository, cache cache.Cache) (*Service, error) {
	s := &Service{
		config: config,
		cache:  cache,
	}

	if err := s.warmConfigCache(); err != nil {
		return nil, err
	}

	return s, nil
}
