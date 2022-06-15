package main

import (
	"errors"
	"os"
)

func readEnv() []error {
	var errs []error

	mysqlUserKey := "MYSQL_USER"
	if os.Getenv(mysqlUserKey) == ""{
		errs = append(errs, errors.New(mysqlUserKey))
	}

	mysqlPasswordKey := "MYSQL_PASSWORD"
	if os.Getenv(mysqlPasswordKey) == ""{
		errs = append(errs, errors.New(mysqlPasswordKey))
	}

	mysqlPortKey := "MYSQL_PORT"
	if os.Getenv(mysqlPortKey) == ""{
		errs = append(errs, errors.New(mysqlPortKey))
	}

	mysqlHostKey := "MYSQL_HOST"
	if os.Getenv(mysqlHostKey) == ""{
		errs = append(errs, errors.New(mysqlHostKey))
	}

	dbNameKey := "DB_NAME"
	if os.Getenv(dbNameKey) == ""{
		errs = append(errs, errors.New(dbNameKey))
	}

	return errs
}
