package models

// User : an application user
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GHUser : a GitHub user
type GHUser struct {
	Login    string `json:"login"`
	Image    string `json:"avatar_url"`
	Blog     string `json:"blog"`
	Location string `json:"location"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
}
