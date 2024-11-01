// services/transactionCreationService.go
package services

import (
	"errors"
	"go-server/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewTransaction function creates a transaction and checks if the sender and receiver are different.
func NewTransaction(sender, receiver primitive.ObjectID, amount float64) (*models.Transaction, error) {
	if sender == receiver {
		return nil, errors.New("sender and receiver cannot be the same user")
	}

	return &models.Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Timestamp: time.Now(),
		Status:    "pending", // Default status set to 'pending'
	}, nil
}
