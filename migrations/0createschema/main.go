package main

import (
	"log"

	"github.com/socialement-competents/goauth/database"
)

// HandleRequest : handle the incoming requests
func main() {
	connection, err := database.Connect()
	if err != nil {
		log.Println("connecting to the database failed: ", err)
		return
	}

	defer connection.Close()

	query := `
		CREATE TABLE IF NOT EXISTS Users (
			id SERIAL PRIMARY KEY NOT NULL,
			email VARCHAR (255) UNIQUE NOT NULL,
			login VARCHAR (255) NOT NULL,
			created TIMESTAMP,
			last_login TIMESTAMP
		)
	`

	_, err = connection.Exec(query)

	if err != nil {
		log.Print("creating the table failed: ", err)
	}
}
