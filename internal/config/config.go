package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBConfig         DBConfig
	AWSRegion        string `env:"AWS_REGION,required"`
	AWSDefaultConfig aws.Config
}

type DBConfig struct {
	Host     string `env:"DB_HOST,required"`
	Name     string `env:"DB_NAME,required"`
	UserName string `env:"DB_USER,required"`
	Port     string `env:"DB_PORT,required" envDefault:"3306"`
}

func (d *DBConfig) Endpoint() string {
	return fmt.Sprintf("%s:%s", d.Host, d.Port)
}

func NewConfig(ctx context.Context) (*Config, error) {
	var cfg Config
	if cerr := env.Parse(&cfg); cerr != nil {
		return nil, cerr
	}

	awscfg, lerr := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.AWSRegion))
	if lerr != nil {
		return nil, lerr
	}

	cfg.AWSDefaultConfig = awscfg

	return &cfg, nil
}
