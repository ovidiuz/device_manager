package interfaces

import (
	"context"

	"github.com/ovidiuz/device_manager/domain"
)

type UserRepo interface {
	GetUser(ctx context.Context, userId string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	SaveUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, email string) error
}
