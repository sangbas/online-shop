package entity

// User represents a user.
type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	FullName string `db:"fullname"`
	Phone    string `db:"phone"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Token    string `db:"token"`
}

// GetID returns the user ID.
func (u User) GetID() string {
	return u.ID
}

// GetName returns the user name.
func (u User) GetUsername() string {
	return u.Username
}
