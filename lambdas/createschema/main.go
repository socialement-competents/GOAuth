package createschema

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	// populate the sql package with pgsql
	_ "github.com/lib/pq"

	"github.com/aws/aws-lambda-go/lambda"
)

// HandleRequest : handle the incoming requests
func HandleRequest(ctx context.Context) error {
	host := os.Getenv("DATABASE_HOST")
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	db := os.Getenv("DATABASE_DATABASE")

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		username, password, db)
	connection, err := sql.Open(host, dbinfo)
	if err != nil {
		log.Print("connecting to the database failed: ", err)
		return err
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

	return err
}

func main() {
	lambda.Start(HandleRequest)
}
