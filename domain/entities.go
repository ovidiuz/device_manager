package domain

type Device struct {
	Id   string `json:"id" db:"id"`
	User string `json:"user" db:"user"`
}

type User struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
