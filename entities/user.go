package entities

import "github.com/fanky5g/ponzu/tokens"

type User struct {
	ID    string `json:"-"`
	Email string `json:"email"`
}

func (*User) EntityName() string {
	return "User"
}

func (*User) GetRepositoryToken() tokens.RepositoryToken {
	return tokens.UserRepositoryToken
}

func (user *User) ItemID() string {
	return user.ID
}

func (user *User) SetItemID(id string) {
	user.ID = id
}
