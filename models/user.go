package models

import "time"

// User : an application user
type User struct {
	ID        int       `json:"id" db:"id"`
	LastLogin time.Time `json:"last_login" db:"last_login"`
	Created   time.Time `json:"created" db:"created"`
	*GHUser
}

// GHUser : a GitHub user
type GHUser struct {
	Bio      string `json:"bio" db:"bio"`
	Blog     string `json:"blog" db:"blog"`
	Email    string `json:"email" db:"email"`
	Image    string `json:"avatar_url" db:"image"`
	Location string `json:"location" db:"location"`
	Login    string `json:"login" db:"login"`
	Name     string `json:"name" db:"name"`
}
