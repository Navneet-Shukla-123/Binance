package models

import "time"

// User model for Email and Password for user
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GetResponse model will hold response from the url
type GetResponse struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// Tokendata model for the token of the user
type Tokendata struct {
	Email      string    `json:"email"`
	Token      string    `json:"token"`
	Expires_At time.Time `json:"expires_at"`
}

type LoginUser struct {
	LoginMethod int    `json:"login_method"` // 1-token,2-cookie
	Email       string `json:"email"`
	Password    string `json:"password"`
}
