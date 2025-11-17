package db

import (
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

var DB *sql.DB

func Init() *sql.DB {
	var err error
	DB, err = sql.Open("sqlite3", "./rookie.db")
	if err != nil {
		panic(err)
	}

	if _, err := DB.Exec(schema); err != nil {
		panic(err)
	}

	return DB
}
