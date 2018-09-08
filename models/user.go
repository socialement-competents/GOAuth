package models

import "time"

// User : an application user
type User struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	LastLogin time.Time `json:"last_login" db:"last_login"`
	Created   time.Time `json:"created" db:"created"`
	*GHUser
}

// GHUser : a GitHub user
type GHUser struct {
	Login    string `json:"login" db:"login"`
	Image    string `json:"avatar_url" db:"image"`
	Blog     string `json:"blog" db:"blog"`
	Location string `json:"location" db:"location"`
	Email    string `json:"email" db:"email"`
	Bio      string `json:"bio" db:"bio"`
}
