package auth

import (
	"github.com/fanky5g/ponzu/internal/database"
)

type UserService struct {
	repository database.Repository
}

func (s *UserService) CreateUser(email string) (*User, error) {
	user, err := s.repository.Insert(&User{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return user.(*User), nil
}

func (s *UserService) DeleteUser(email string) error {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	return s.repository.DeleteById(user.ID)
}

func (s *UserService) UpdateUser(user, update *User) error {
	_, err := s.repository.UpdateById(user.ID, update)
	return err
}

func (s *UserService) GetUserByEmail(email string) (*User, error) {
	u, err := s.repository.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, nil
	}

	return u.(*User), nil
}

func (s *UserService) ListUsers() ([]User, error) {
	uu, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	users := make([]User, len(uu))
	for i := range uu {
		u := uu[i].(*User)
		users[i] = *u
	}

	return users, nil
}

func NewUserService(db database.Database) (*UserService, error) {
	return &UserService{repository: db.GetRepositoryByToken(UserRepositoryToken)}, nil
}
