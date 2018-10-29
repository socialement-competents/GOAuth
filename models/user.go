package models

import (
	"fmt"
	"time"
)

// User : an application user
type User struct {
	ID        int       `json:"id"`
	LastLogin time.Time `json:"last_login"`
	Created   time.Time `json:"created"`
	Provider  string    `json:"provider"`
	*GHUser
	*FitBitUser
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

// FitBitUser : a Fitbit user
type FitBitUser struct {
	Age        int    `json:"age" db:"fitbit_age"`
	Avatar     string `json:"avatar150" db:"fitbit_avatar150"`
	FullName   string `json:"fullName" db:"fitbit_fullname"`
	EncodedID  string `json:"encodedId" db:"fitbit_id"`
	RawPayload string `json:"raw_payload" db:"fitbit_json_payload"`
}

// GetUniqueIdentifier returns the unique of the user no matter the provider (email / ID / login / ...)
func (u *User) GetUniqueIdentifier() (string, error) {
	switch {
	case u.Provider == GithubProvider:
		return u.GHUser.Email, nil
	case u.Provider == FitBitProvider:
		return u.FitBitUser.EncodedID, nil
	default:
		return "", fmt.Errorf("No unique identifier for user %d / provider %s", u.ID, u.Provider)
	}
}

// GetUniqueIdentifierName returns the column name for the unique identifier
func (u *User) GetUniqueIdentifierName() (string, error) {
	switch {
	case u.Provider == GithubProvider:
		return "login", nil
	case u.Provider == FitBitProvider:
		return "fitbit_id", nil
	default:
		return "", fmt.Errorf("No unique identifier provider %s", u.Provider)
	}
}

// GetImage returns the image of the user no matter the provider
func (u *User) GetImage() string {
	switch {
	case u.GHUser.Image != "":
		return u.GHUser.Image
	case u.FitBitUser.Avatar != "":
		return u.FitBitUser.Avatar
	default:
		return ""
	}
}

// GetName returns the name of the user no matter the provider
func (u *User) GetName() string {
	switch {
	case u.GHUser.Name != "":
		return u.GHUser.Name
	case u.FitBitUser.FullName != "":
		return u.FitBitUser.FullName
	default:
		return ""
	}
}
