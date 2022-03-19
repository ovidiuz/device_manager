package domain

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type BasicAuth struct {
	Email    string `json:"email" validate:"required" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"11111111"`
}

type RegisterRequest struct {
	BasicAuth
}

type LoginRequest struct {
	BasicAuth
}

func (a *BasicAuth) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(a.Email, is.Email),
		validation.Field(a.Password, validation.Length(minPasswordLength, 0)),
	)
}
