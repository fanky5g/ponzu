package repositories

import (
	"github.com/fanky5g/ponzu/entities"
)

type UserRepositoryInterface interface {
	SetUser(usr *entities.User) error
	UpdateUser(usr, updatedUsr *entities.User) error
	DeleteUser(email string) error
	GetUserByEmail(email string) (*entities.User, error)
	GetAllUsers() ([][]byte, error)
}
