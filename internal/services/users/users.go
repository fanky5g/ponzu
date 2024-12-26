package users

import (
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/internal/auth"
	"github.com/fanky5g/ponzu/tokens"
)

type service struct {
	repository driver.Repository
}

type Service interface {
	CreateUser(email string) (*auth.User, error)
	DeleteUser(email string) error
	UpdateUser(user, update *auth.User) error
	ListUsers() ([]auth.User, error)
	GetUserByEmail(email string) (*auth.User, error)
}

func (s *service) CreateUser(email string) (*auth.User, error) {
	user, err := s.repository.Insert(&auth.User{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return user.(*auth.User), nil
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

func (s *service) UpdateUser(user, update *auth.User) error {
	_, err := s.repository.UpdateById(user.ID, update)
	return err
}

func (s *service) GetUserByEmail(email string) (*auth.User, error) {
	u, err := s.repository.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, nil
	}

	return u.(*auth.User), nil
}

func (s *service) ListUsers() ([]auth.User, error) {
	uu, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	users := make([]auth.User, len(uu))
	for i := range uu {
		if err != nil {
			return nil, err
		}

		u := uu[i].(*auth.User)
		users[i] = *u
	}

	return users, nil
}

func New(db driver.Database) (Service, error) {
	return &service{repository: db.GetRepositoryByToken(tokens.UserRepositoryToken)}, nil
}
