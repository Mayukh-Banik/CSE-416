package models

import (
	"time"
)

// User defines the structure for user data
type User struct {
	UUID         string            `json:"uuid"`         // Unique internal identifier
	Address      string            `json:"address"`      // Wallet address
	PublicKey    string            `json:"public_key"`   // Public key (optional)
	CreatedDate  time.Time         `json:"created_date"` // Timestamp of user creation
	Metadata     map[string]string `json:"metadata"`     // Additional user-specific data
	Balance      float64           `json:"balance"`      // Optional, tracks wallet balance
	Transactions []string          `json:"transactions"` // List of associated transaction IDs
}
