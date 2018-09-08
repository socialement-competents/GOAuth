package database

import (
	"errors"
	"time"

	"github.com/socialement-competents/goauth/models"
)

// CreateUser inserts a new User in the database
func (c *Client) CreateUser(u *models.User) (int, error) {
	query := `
		INSERT INTO users (bio, blog, email, image, location, login, name, provider, last_login, created)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id;
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
		u.Provider,
		u.LastLogin,
		u.Created,
	).Scan(&id)
	return id, err
}

// UpdateUser will update a row from its ID
func (c *Client) UpdateUser(u *models.User) error {
	query := `
		UPDATE users
		SET	
			bio = $2,
			blog = $3,
			email = $4,
			image = $5,
			location = $6,
			login = $7,
			name = $8,
			last_login = $9
		WHERE id = $1;
	`

	_, err := c.Connection.Exec(
		query,
		u.ID,
		u.Bio,
		u.Blog,
		u.Email,
		u.Image,
		u.Location,
		u.Login,
		u.Name,
		u.LastLogin,
	)

	return err
}

// GetUserByLogin selects an user from his login
func (c *Client) GetUserByLogin(login, provider string) (*models.User, error) {
	query := `
		SELECT id, bio, blog, email, image, location, login, name, provider, last_login, created
		FROM users
		WHERE login = $1 AND provider = $2;
	`
	u := models.User{GHUser: &models.GHUser{}}
	row := c.Connection.QueryRow(query, login, provider)
	if row == nil {
		return nil, errors.New("Not found")
	}

	if err := row.Scan(
		&u.ID,
		&u.Bio,
		&u.Blog,
		&u.Email,
		&u.Image,
		&u.Location,
		&u.Login,
		&u.Name,
		&u.Provider,
		&u.LastLogin,
		&u.Created,
	); err != nil {
		return nil, err
	}

	return &u, nil
}
