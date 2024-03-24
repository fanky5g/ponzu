package users

import (
	"encoding/json"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
	"github.com/fanky5g/ponzu/tokens"
)

type service struct {
	repository repositories.UserRepositoryInterface
}

type Service interface {
	CreateUser(email string) (*entities.User, error)
	DeleteUser(email string) error
	UpdateUser(user, update *entities.User) error
	ListUsers() ([]entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
}

func (s *service) CreateUser(email string) (*entities.User, error) {
	user := &entities.User{
		Email: email,
	}

	err := s.repository.SetUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) DeleteUser(email string) error {
	return s.repository.DeleteUser(email)
}

func (s *service) UpdateUser(user, update *entities.User) error {
	return s.repository.UpdateUser(user, update)
}

func (s *service) GetUserByEmail(email string) (*entities.User, error) {
	return s.repository.GetUserByEmail(email)
}

func (s *service) ListUsers() ([]entities.User, error) {
	// get all users to list
	jj, err := s.repository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var users []entities.User
	for i := range jj {
		var u entities.User
		err = json.Unmarshal(jj[i], &u)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func New(db driver.Database) (Service, error) {
	return &service{repository: db.Get(tokens.UserRepositoryToken).(repositories.UserRepositoryInterface)}, nil
}
