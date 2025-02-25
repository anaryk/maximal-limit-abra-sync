package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Connector struct {
	db *sql.DB
}

func NewMySQLConnector(dbname, dbhost, dbuser, dbpassword string) (*Connector, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbuser, dbpassword, dbhost, dbname))
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &Connector{db: db}, nil
}
