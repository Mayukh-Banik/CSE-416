package models

import "time"

type Challenge struct {
    Address     string    `json:"address"`      // Wallet address of the user
    Challenge   string    `json:"challenge"`    // Random string to sign
    Expiry      time.Time `json:"expiry"`       // Expiry time for the challenge
}