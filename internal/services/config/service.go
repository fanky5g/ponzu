package config

import (
	"fmt"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/models"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/fanky5g/ponzu/util"
	"reflect"
)

type service struct {
	repository driver.Repository
}

type Service interface {
	GetAppName() (string, error)
	SetConfig(config *entities.Config) error
	Get() (*entities.Config, error)
}

func (s *service) GetStringValue(key string) (string, error) {
	cfg, err := s.repository.Latest()
	if err != nil {
		return "", err
	}

	return util.StringFieldByJSONTagName(cfg, key)
}

func (s *service) GetBoolValue(key string) (bool, error) {
	cfg, err := s.repository.Latest()
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

func (s *service) SetConfig(config *entities.Config) error {
	cfg, err := s.Get()
	if err != nil {
		return err
	}

	if cfg == nil {
		_, err = s.repository.Insert(config)
		return err
	}

	_, err = s.repository.UpdateById(cfg.ID, config)
	return err
}

func (s *service) Get() (*entities.Config, error) {
	cfg, err := s.repository.Latest()
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		return nil, nil
	}

	return cfg.(*entities.Config), nil
}

func (s *service) GetAppName() (string, error) {
	return s.GetStringValue("name")
}

func New(db driver.Database) (Service, error) {
	return &service{repository: db.Get(
		models.WrapPonzuModelNameSpace(tokens.ConfigRepositoryToken),
	).(driver.Repository)}, nil
}
