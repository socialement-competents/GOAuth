package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	// populate the sql package with pgsql
	_ "github.com/lib/pq"
)

// Connect : Reads the environment variables to return a SQL db connection
func Connect() (*sql.DB, error) {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	db := os.Getenv("DATABASE_DATABASE")

	if host == "" {
		return nil, errors.New("$DATABASE_HOST should be set")
	}
	if username == "" {
		return nil, errors.New("$DATABASE_USERNAME should be set")
	}
	if password == "" {
		return nil, errors.New("$DATABASE_PASSWORD should be set")
	}
	if db == "" {
		return nil, errors.New("$DATABASE_DATABASE should be set")
	}

	dbinfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		username,
		password,
		db,
	)
	return sql.Open("postgres", dbinfo)
}
