package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// populate the sql package with pgsql
	_ "github.com/lib/pq"
)

// Client : SQL Client
type Client struct {
	Connection *sql.DB
}

// NewClient : Reads the environment variables to return a SQL client
func NewClient() (*Client, error) {
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
	connection, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}

	client := Client{Connection: connection}
	client.cleanupHook()

	return &client, nil
}

// Intercept the exit signal and close the connection properly before exiting
func (client *Client) cleanupHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		client.Connection.Close()
		os.Exit(0)
	}()
}
