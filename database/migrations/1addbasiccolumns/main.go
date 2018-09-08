package main

import (
	"log"

	"github.com/socialement-competents/goauth/database"
)

// Add some useful fields to the user table
func main() {
	client, err := database.NewClient()
	if err != nil {
		log.Println("connecting to the database failed: ", err)
		return
	}

	query := `
		ALTER TABLE Users
		ADD COLUMN image VARCHAR(500),
		ADD COLUMN blog VARCHAR(500),
		ADD COLUMN location VARCHAR(500),
		ADD COLUMN bio TEXT;
	`

	_, err = client.Connection.Exec(query)

	if err != nil {
		log.Print("creating the table failed: ", err)
	}
}
