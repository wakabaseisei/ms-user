package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/caarlos0/env/v11"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
)

func handler(ctx context.Context, event json.RawMessage) error {
	var cfg Config
	if cerr := env.Parse(&cfg); cerr != nil {
		return fmt.Errorf("parse env: %v", cerr)
	}

	token, gerr := generateAuthToken(cfg.DBConfig.Endpoint(), cfg.DBConfig.UserName, cfg.DBConfig.Region)
	if gerr != nil {
		return fmt.Errorf("generate IAM auth token: %v", gerr)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&multiStatements=true&allowCleartextPasswords=true",
		cfg.DBConfig.UserName, token, cfg.DBConfig.Endpoint(), cfg.DBConfig.Name)

	db, serr := sql.Open("mysql", dsn)
	if serr != nil {
		return fmt.Errorf("connect to DB: %v", serr)
	}
	defer db.Close()

	if merr := runMigration(db, cfg.DBConfig.Name); merr != nil {
		return fmt.Errorf("migration: %v", merr)
	}

	log.Println("Migration was successful!")
	return nil
}

func generateAuthToken(host, user, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "", err
	}

	authenticationToken, terr := auth.BuildAuthToken(
		context.TODO(), host, region, user, cfg.Credentials)
	if terr != nil {
		return "", terr
	}

	return authenticationToken, nil
}

func runMigration(db *sql.DB, dbName string) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}
	migrationDir := filepath.Join(filepath.Dir(exePath), "../db/migrations")

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		dbName,
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
