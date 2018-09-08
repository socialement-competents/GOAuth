package database

import (
	"time"

	"github.com/socialement-competents/goauth/models"
)

// Create inserts a new User in the database
func (c *Client) Create(u *models.User) (int, error) {
	query := `
		INSERT INTO Users (bio, blog, email, image, location, login, name, last_login, created)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	u.Created = time.Now()

	id := 0
	err := c.Connection.QueryRow(
		query,
		u.Bio,
		u.Blog,
		u.Email,
		u.Image,
		u.Location,
		u.Login,
		u.Name,
		u.LastLogin,
		u.Created,
	).Scan(&id)
	return id, err
}

// GetUserByLogin selects an user from his login
func (c *Client) GetUserByLogin(login string) (*models.User, error) {
	query := `
		SELECT id, bio, blog, email, image, location, login, name, last_login, created
		FROM User
		WHERE login = $1
	`
	u := models.User{}
	err := c.Connection.QueryRow(query, login).Scan(
		&u.ID,
		&u.Bio,
		&u.Blog,
		&u.Email,
		&u.Image,
		&u.Location,
		&u.Login,
		&u.Name,
		&u.LastLogin,
		&u.Created,
	)
	return &u, err
}
