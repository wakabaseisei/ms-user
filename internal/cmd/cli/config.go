package main

import "fmt"

type Config struct {
	DBConfig DBConfig
}

type DBConfig struct {
	Host     string `env:"DB_HOST,required"`
	Name     string `env:"DB_NAME,required"`
	UserName string `env:"DB_USER,required"`
	Region   string `env:"AWS_REGION,required"`
	Port     string `env:"DB_PORT,required" envDefault:"3306"`
}

func (d *DBConfig) Endpoint() string {
	return fmt.Sprintf("%s:%s", d.Host, d.Port)
}
