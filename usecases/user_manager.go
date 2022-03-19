package usecases

import (
	"context"

	"github.com/ovidiuz/device_manager/domain"
	"github.com/ovidiuz/device_manager/gateways/interfaces"
)

type UserManager struct {
	userRepo interfaces.UserRepo
}

func NewUserManager(userRepo interfaces.UserRepo) *UserManager {
	return &UserManager{
		userRepo: userRepo,
	}
}

func (m *UserManager) GetUser(ctx context.Context, userId string) (*domain.User, error) {
	return m.userRepo.GetUser(ctx, userId)
}

func (m *UserManager) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.userRepo.GetUserByEmail(ctx, email)
}
