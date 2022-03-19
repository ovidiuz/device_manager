package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type BasicAuth struct {
	Email    string `json:"email" validate:"required" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"11111111"`
}

type RegisterRequest struct {
	BasicAuth
	PasswordConfirm string `json:"password_confirm" validate:"required" example:"11111111"`
}

type LoginRequest struct {
	BasicAuth
}

func (a *BasicAuth) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Email, is.EmailFormat.Error("email must be a valid email address")),
		//validation.Field(&a.Email, is.Email.Error("email must be an existing email address")),
		validation.Field(&a.Password, validation.Length(minPasswordLength, 0).Error("password must have at least 8 characters")),
	)
}

func (r *RegisterRequest) Validate() error {
	if err := r.BasicAuth.Validate(); err != nil {
		return err
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.PasswordConfirm, validation.In(r.Password).Error("passwords do not match")),
	)
}
