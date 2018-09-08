package models

import "time"

// User : an application user
type User struct {
	ID        int       `json:"id"`
	LastLogin time.Time `json:"last_login"`
	Created   time.Time `json:"created"`
	Provider  string    `json:"provider"`
	*GHUser
}

// GHUser : a GitHub user
type GHUser struct {
	Bio      string `json:"bio"`
	Blog     string `json:"blog"`
	Email    string `json:"email"`
	Image    string `json:"avatar_url" db:"image"`
	Location string `json:"location"`
	Login    string `json:"login"`
	Name     string `json:"name"`
}
