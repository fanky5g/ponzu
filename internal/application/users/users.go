package users

import (
	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/entities"
)

type Service struct {
	users database.Repository
}

func (s *Service) CreateUser(email string) (*entities.User, error) {
	user, err := s.users.Insert(&entities.User{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return user.(*entities.User), nil
}

func (s *Service) DeleteUser(email string) error {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	return s.users.DeleteById(user.ID)
}

func (s *Service) UpdateUser(user, update *entities.User) error {
	_, err := s.users.UpdateById(user.ID, update)
	return err
}

func (s *Service) GetUserByEmail(email string) (*entities.User, error) {
	u, err := s.users.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, nil
	}

	return u.(*entities.User), nil
}

func (s *Service) ListUsers() ([]entities.User, error) {
	uu, err := s.users.FindAll()
	if err != nil {
		return nil, err
	}

	users := make([]entities.User, len(uu))
	for i := range uu {
		if err != nil {
			return nil, err
		}

		u := uu[i].(*entities.User)
		users[i] = *u
	}

	return users, nil
}

func New(users database.Repository) (*Service, error) {
	return &Service{users: users}, nil
}
