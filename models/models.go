package models

import "time"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetResponse struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type Tokendata struct {
	Email      string    `json:"email"`
	Token      string    `json:"token"`
	Expires_At time.Time `json:"expires_at"`
}
