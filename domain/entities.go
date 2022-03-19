package domain

type Device struct {
	DeviceID string `json:"device_id" db:"device_id"`
	UserID   string `json:"user_id" db:"user_id"`
}

type User struct {
	UserID   string `json:"user_id" db:"user_id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
}
