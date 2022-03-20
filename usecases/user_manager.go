package usecases

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/casbin/casbin/v2"

	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"

	"github.com/ovidiuz/device_manager/domain"
	"github.com/ovidiuz/device_manager/gateways/interfaces"
)

const secretKey = "secret"

type UserManager struct {
	userRepo    interfaces.UserRepo
	enforcer    casbin.IEnforcer
	jwtTokenTTL time.Duration
}

func NewUserManager(userRepo interfaces.UserRepo, enforcer casbin.IEnforcer) *UserManager {
	return &UserManager{
		userRepo: userRepo,
		enforcer: enforcer,
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
	if _, err := m.enforcer.AddGroupingPolicy(user.UserID, user.Role); err != nil {
		logger.WithError(err).Errorf("could not add role inheritance for user=%s role=%s", user.UserID, user.Role)
		return nil, err
	}

	logger.Infof("user with email=%s successfully registered", user.Email)
	return user, nil
}

func (m *UserManager) LoginUser(ctx context.Context, request *domain.LoginRequest) (string, error) {
	logger := logrus.WithContext(ctx)
	user, err := m.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		logger.WithError(err).Error("could not get user from DB")
		return "", err
	}

	// Compare the passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		logger.WithError(err).Debug("wrong password")
		return "", err
	}

	// Generate the JWT token
	token, err := generateJWT(user.UserID, m.jwtTokenTTL)
	if err != nil {
		logger.WithError(err).Errorf("could not generate JWT for user=%s", user.UserID)
		return "", err
	}

	logger.Debugf("successfully generated token for user=%s", user.UserID)
	return token, nil
}

func generateJWT(userID string, ttl time.Duration) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Issuer:    userID,
	})

	return claims.SignedString([]byte(secretKey))
}
