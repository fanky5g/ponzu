package users

import (
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
	"github.com/fanky5g/ponzu/tokens"
)

type service struct {
	repository repositories.GenericRepositoryInterface
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

	_, err := s.repository.Insert(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) DeleteUser(email string) error {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	return s.repository.DeleteById(user.ID)
}

func (s *service) UpdateUser(user, update *entities.User) error {
	_, err := s.repository.UpdateById(user.ID, update)
	return err
}

func (s *service) GetUserByEmail(email string) (*entities.User, error) {
	u, err := s.repository.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, nil
	}

	return u.(*entities.User), nil
}

func (s *service) ListUsers() ([]entities.User, error) {
	uu, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	users := make([]entities.User, 0)
	for i := range uu {
		if err != nil {
			return nil, err
		}

		u := uu[i].(*entities.User)
		users[i] = *u
	}

	return users, nil
}

func New(db driver.Database) (Service, error) {
	return &service{repository: db.Get(tokens.UserRepositoryToken).(repositories.GenericRepositoryInterface)}, nil
}
