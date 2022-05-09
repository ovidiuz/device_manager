package repositories

import (
	"context"
	"database/sql"

	"github.com/ovidiuz/device_manager/device_manager/domain"
	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
)

// TODO: add metrics and timers for Prometheus

const (
	getUserStmt        = "SELECT * FROM users WHERE user_id=$1"
	getUserByEmailStmt = "SELECT * FROM users WHERE email=$1"
	saveUserStmt       = "INSERT INTO users (email, role, password) VALUES (:email, :role, :password)"
	deleteUserStmt     = "DELETE FROM users WHERE email=$1"
)

type UserSQLRepo struct {
	db *sqlx.DB
}

func NewUserSQLRepo(db *sqlx.DB) *UserSQLRepo {
	return &UserSQLRepo{
		db: db,
	}
}

func (r *UserSQLRepo) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.GetContext(ctx, user, getUserStmt, userID)
	if err == sql.ErrNoRows {
		logrus.WithContext(ctx).Debugf("user=%s does not exist", userID)
		return nil, domain.ErrNotFound
	} else if err != nil {
		logrus.WithContext(ctx).WithError(err).Errorf("could not get user=%s", userID)
		return nil, err
	}
	return user, nil
}

func (r *UserSQLRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.GetContext(ctx, user, getUserByEmailStmt, email)
	if err == sql.ErrNoRows {
		logrus.WithContext(ctx).Debugf("user with email=%s does not exist", email)
		return nil, domain.ErrNotFound
	} else if err != nil {
		logrus.WithContext(ctx).WithError(err).Errorf("could not get user with email=%s", email)
		return nil, err
	}
	return user, nil
}

func (r *UserSQLRepo) SaveUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.NamedExecContext(ctx, saveUserStmt, user)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Errorf("could not save user with email=%s", user.Email)
		return err
	}
	return nil
}

func (r *UserSQLRepo) DeleteUser(ctx context.Context, email string) error {
	_, err := r.db.ExecContext(ctx, deleteUserStmt, email)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Errorf("could not delete user with email=%s", email)
		return err
	}
	return nil
}
