package entities

type User struct {
	ID    string `json:"-"`
	Email string `json:"email"`
}

func (*User) EntityName() string {
	return "User"
}

func (*User) GetRepositoryToken() string {
	return "users"
}

func (user *User) ItemID() string {
	return user.ID
}

func (user *User) SetItemID(id string) {
	user.ID = id
}
