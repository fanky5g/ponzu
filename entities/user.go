package entities

import "github.com/fanky5g/ponzu/tokens"

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (*User) GetRepositoryToken() tokens.RepositoryToken {
	return tokens.UserRepositoryToken
}
