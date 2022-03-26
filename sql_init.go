package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ovidiuz/device_manager/domain"
)

const (
	defaultSQLTimeout       = 5 * time.Second
	createExtensionUUIDStmt = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	createUsersTableStmt    = `CREATE TABLE IF NOT EXISTS users (
    							user_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    							email VARCHAR(255) UNIQUE NOT NULL,
    							password VARCHAR(255) NOT NULL,
    							role VARCHAR(255) NOT NULL
                   			);`
	createDevicesTableStmt = `CREATE TABLE IF NOT EXISTS devices (
								device_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
								user_id uuid NOT NULL REFERENCES users(user_id)
							);`
)

func initSQL(conf *domain.ServiceConfig) (*sqlx.DB, error) {
	connInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", conf.PSQLHost, conf.PSQLUsername, conf.PSQLPassword, conf.PSQLDBName)
	dbConn, err := sqlx.Open("postgres", connInfo)
	if err != nil {
		logrus.WithError(err).Errorf("could not connect to the PostgreSQL database %s", conf.PSQLDBName)
		return nil, err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultSQLTimeout)
	defer cancelFunc()

	err = installUUIDExtension(ctx, dbConn)
	if err != nil {
		logrus.WithError(err).Error("could not install uuid-ossp extension")
		return nil, err
	}
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

func installUUIDExtension(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, createExtensionUUIDStmt)
	return err
}

func createUsersTable(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, createUsersTableStmt)
	return err
}

func createDevicesTable(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, createDevicesTableStmt)
	return err
}
