package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/go-sql-driver/mysql"

	"github.com/wakabaseisei/ms-user/internal/config"
)

func NewDatabase(ctx context.Context, cfg config.DBConfig, awscfg aws.Config) (*sql.DB, error) {
	token, terr := generateAuthToken(ctx, cfg.Endpoint(), cfg.UserName, awscfg)
	if terr != nil {
		return nil, fmt.Errorf("generate IAM Auth Token: %v", terr)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&multiStatements=true&allowCleartextPasswords=true&parseTime=true",
		cfg.UserName, token, cfg.Endpoint(), cfg.Name)

	db, serr := sql.Open("mysql", dsn)
	if serr != nil {
		return nil, fmt.Errorf("connect to DB: %v", serr)
	}

	return db, nil
}

func generateAuthToken(ctx context.Context, endpoint, user string, awscfg aws.Config) (string, error) {
	authenticationToken, terr := auth.BuildAuthToken(
		ctx, endpoint, awscfg.Region, user, awscfg.Credentials)
	if terr != nil {
		return "", terr
	}

	return authenticationToken, nil
}
