package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) CloseDB() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
