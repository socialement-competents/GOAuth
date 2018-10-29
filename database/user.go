package database

import (
	"errors"
	"time"

	"github.com/socialement-competents/goauth/models"
)

// CreateUser inserts a new User in the database
func (c *Client) CreateUser(u *models.User) (int, error) {
	query := `
		INSERT INTO users (bio, blog, email, image, location, login, name, fitbit_age, fitbit_avatar150, fitbit_id, fitbit_fullname, fitbit_json_payload, provider, last_login, created)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id;
	`

	u.Created = time.Now()

	id := 0
	err := c.Connection.QueryRow(
		query,
		u.GHUser.Bio,
		u.GHUser.Blog,
		u.GHUser.Email,
		u.GHUser.Image,
		u.GHUser.Location,
		u.GHUser.Login,
		u.GHUser.Name,
		u.FitBitUser.Age,
		u.FitBitUser.Avatar,
		u.FitBitUser.EncodedID,
		u.FitBitUser.FullName,
		u.FitBitUser.RawPayload,
		u.Provider,
		u.LastLogin,
		u.Created,
	).Scan(&id)
	u.ID = id
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
			fitbit_age = $9,
			fitbit_avatar150 = $10,
			fitbit_id = $11,
			fitbit_fullname = $12,
			fitbit_json_payload = $13,
			last_login = $14
		WHERE id = $1;
	`

	_, err := c.Connection.Exec(
		query,
		u.ID,
		u.GHUser.Bio,
		u.GHUser.Blog,
		u.GHUser.Email,
		u.GHUser.Image,
		u.GHUser.Location,
		u.GHUser.Login,
		u.GHUser.Name,
		u.FitBitUser.Age,
		u.FitBitUser.Avatar,
		u.FitBitUser.EncodedID,
		u.FitBitUser.FullName,
		u.FitBitUser.RawPayload,
		u.LastLogin,
	)

	return err
}

// GetUserByIdentifier selects an user from his FitBit ID
func (c *Client) GetUserByIdentifier(column, id string) (*models.User, error) {
	query := `
		SELECT id, provider, last_login, created, fitbit_age, fitbit_avatar150, fitbit_fullname, fitbit_id, fitbit_json_payload, bio, blog, email, image, location, login, name
		FROM users
		WHERE $1 = $2;
	`
	u := models.User{
		GHUser:     &models.GHUser{},
		FitBitUser: &models.FitBitUser{},
	}
	row := c.Connection.QueryRow(query, column, id)
	if row == nil {
		return nil, errors.New("Not found")
	}

	err := row.Scan(
		&u.ID,
		&u.Provider,
		&u.LastLogin,
		&u.Created,
		&u.FitBitUser.Age,
		&u.FitBitUser.Avatar,
		&u.FitBitUser.FullName,
		&u.FitBitUser.EncodedID,
		&u.FitBitUser.RawPayload,
		&u.GHUser.Bio,
		&u.GHUser.Blog,
		&u.GHUser.Email,
		&u.GHUser.Image,
		&u.GHUser.Location,
		&u.GHUser.Login,
		&u.GHUser.Name,
	)

	return &u, err
}

// GetUserByLogin selects an user from his login
func (c *Client) GetUserByLogin(login, provider string) (*models.User, error) {
	query := `
		SELECT id, bio, blog, email, image, location, login, name, provider, last_login, created
		FROM users
		WHERE login = $1 AND provider = $2;
	`
	u := models.User{
		GHUser:     &models.GHUser{},
		FitBitUser: &models.FitBitUser{},
	}
	row := c.Connection.QueryRow(query, login, provider)
	if row == nil {
		return nil, errors.New("Not found")
	}

	err := row.Scan(
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
	)

	return &u, err
}
