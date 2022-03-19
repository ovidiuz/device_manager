package usecases

import (
	"context"

	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"

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

func (m *UserManager) RegisterUser(ctx context.Context, request *domain.RegisterRequest) (*domain.User, error) {
	logger := logrus.WithContext(ctx)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), domain.PasswordCost)
	if err != nil {
		logger.WithError(err).Error("could not encrypt password")
		return nil, err
	}

	user := &domain.User{
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	if err := m.userRepo.SaveUser(ctx, user); err != nil {
		logger.WithError(err).Errorf("could not save user with email=%s to DB", user.Email)
		return nil, err
	}

	logger.Infof("user with email=%s successfully registered", user.Email)
	return user, nil
}
