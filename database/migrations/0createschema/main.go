package main

import (
	"log"

	"github.com/socialement-competents/goauth/database"
)

// Create a basic User table
func main() {
	client, err := database.NewClient()
	if err != nil {
		log.Println("connecting to the database failed: ", err)
		return
	}

	query := `
		CREATE TABLE IF NOT EXISTS Users (
			id SERIAL PRIMARY KEY NOT NULL,
			email VARCHAR (255) UNIQUE NOT NULL,
			login VARCHAR (255) NOT NULL,
			created TIMESTAMP,
			last_login TIMESTAMP
		)
	`

	_, err = client.Connection.Exec(query)

	if err != nil {
		log.Print("creating the table failed: ", err)
	}
}
