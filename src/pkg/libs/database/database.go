package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"sync"
	"time"
)

var once sync.Once

func GetDbEngine() *sql.DB {
	var database *sql.DB
	once.Do(func() {
		var err error
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
			os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"), os.Getenv("DB_NAME")))
		if err != nil {
			log.Fatal(fmt.Sprintf("Could not connect to MySQL Instance: %v", err))
		}

		maxOpenConnections := 100
		maxIdleConnections := 10
		maxConnectionLifetime := time.Hour

		log.Println(fmt.Sprintf("mysql - settings max connection lifetime to: %v", maxConnectionLifetime))
		log.Println(fmt.Sprintf("mysql - settings max open idle connections to: %v", maxIdleConnections))
		log.Println(fmt.Sprintf("mysql - settings max open connections to: %v", maxOpenConnections))

		db.SetConnMaxLifetime(maxConnectionLifetime)
		db.SetMaxIdleConns(maxIdleConnections)
		db.SetMaxOpenConns(maxOpenConnections)

		database = db
	})

	return database
}
