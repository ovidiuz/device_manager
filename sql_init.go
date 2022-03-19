package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	"github.com/ovidiuz/device_manager/domain"
)

const (
	defaultSQLTimeout    = 5 * time.Second
	createUsersTableStmt = `CREATE TABLE users IF NOT EXISTS (
    							user_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    							email VARCHAR(255) UNIQUE NOT NULL,
    							password VARCHAR(255) NOT NULL
                   			);`
	createDevicesTableStmt = `CREATE TABLE devices IF NOT EXISTS(
								device_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
								user_id VARCHAR(255) NOT NULL REFERENCES users(user_id)
							);`
)

func initSQL(conf *domain.ServiceConfig) (*sqlx.DB, error) {
	connInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", conf.PSQLHost, conf.PSQLUsername, conf.PSQLPassword, conf.PSQLDBName)
	dbConn, err := sqlx.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultSQLTimeout)
	defer cancelFunc()

	err = createUsersTable(ctx, dbConn)
	if err != nil {
		logrus.WithError(err).Error("could not create table 'users'")
		return nil, err
	}
	err = createDevicesTable(ctx, dbConn)
	if err != nil {
		logrus.WithError(err).Error("could not create table 'devices'")
		return nil, err
	}

	return dbConn, nil
}

func createUsersTable(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, createUsersTableStmt)
	return err
}

func createDevicesTable(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, createDevicesTableStmt)
	return err
}
