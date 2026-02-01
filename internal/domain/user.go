package domain

type User struct {
	id           int
	email        string
	passwordHash string
}

func NewUser(email string, passwordHash string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}

	return &User{
		email:        email,
		passwordHash: passwordHash,
	}, nil
}

func FromStorageUser(id int, email string, passwordHash string) *User {
	return &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
	}
}

func (u *User) UserID() int {
	return u.id
}
func (u *User) Email() string {
	return u.email
}
func (u *User) PasswordHash() string {
	return u.passwordHash
}
