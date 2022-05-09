package usecases

import (
	"context"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/ovidiuz/device_manager/device_manager/jwt"

	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"

	"github.com/ovidiuz/device_manager/device_manager/domain"
	"github.com/ovidiuz/device_manager/device_manager/gateways/interfaces"
)

type UserManager struct {
	userRepo    interfaces.UserRepo
	enforcer    casbin.IEnforcer
	jwtTokenTTL time.Duration
}

func NewUserManager(userRepo interfaces.UserRepo, enforcer casbin.IEnforcer, jwtTokenTTL time.Duration) *UserManager {
	return &UserManager{
		userRepo:    userRepo,
		enforcer:    enforcer,
		jwtTokenTTL: jwtTokenTTL,
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
		Role:     domain.UserRole,
		Password: string(hashedPassword),
	}

	if err := m.userRepo.SaveUser(ctx, user); err != nil {
		logger.WithError(err).Errorf("could not save user with email=%s to DB", user.Email)
		return nil, err
	}
	if _, err := m.enforcer.AddGroupingPolicy(user.UserID, string(user.Role)); err != nil {
		logger.WithError(err).Errorf("could not add role inheritance for user=%s role=%s", user.UserID, user.Role)
		return nil, err
	}

	logger.Debugf("user with email=%s successfully registered", user.Email)
	return user, nil
}

func (m *UserManager) LoginUser(ctx context.Context, request *domain.LoginRequest) (string, error) {
	logger := logrus.WithContext(ctx)
	user, err := m.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		logger.WithError(err).Error("could not get user with email=%s from DB", request.Email)
		return "", err
	}

	// Compare the passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		logger.WithError(err).Debug("wrong password")
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return "", domain.ErrUnauthorized
		}
		return "", err
	}

	// Generate the JWT token
	token, err := jwt.GenerateJWT(user.UserID, m.jwtTokenTTL)
	if err != nil {
		logger.WithError(err).Errorf("could not generate JWT for user=%s", user.UserID)
		return "", err
	}

	logger.Debugf("successfully generated JWT token for user=%s", user.UserID)
	return token, nil
}
