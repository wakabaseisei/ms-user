package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	lambda.Start(handler)
}
