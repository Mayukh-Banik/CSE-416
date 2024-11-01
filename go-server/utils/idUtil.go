package utils

import "github.com/google/uuid"

// GenerateUserID generates a unique UUID string
func GenerateUserID() string {
    return uuid.New().String()
}
