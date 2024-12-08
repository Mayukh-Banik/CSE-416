package models

import (
	"time"
)

// User defines the structure for user data
type User struct {
	UUID        string    `json:"uuid"`         // Unique internal identifier
	Address     string    `json:"address"`      // Wallet address
	CreatedDate time.Time `json:"created_date"` // Timestamp of user creation
}
