package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transaction structure represents a single transaction record.
type Transaction struct {
	TransactionID string             `bson:"transaction_id" json:"transaction_id"`
	Sender        primitive.ObjectID `bson:"sender" json:"sender"`
	Receiver      primitive.ObjectID `bson:"receiver" json:"receiver"`
	Amount        float64            `bson:"amount" json:"amount"`
	Timestamp     time.Time          `bson:"timestamp" json:"timestamp"`
	FileName      string             `bson:"file_name" json:"file_name"`
	FileID        string             `bson:"file_id" json:"file_id"`
	FileSize      string             `bson:"file_size" json:"file_size"`
	Fee           float64            `bson:"fee" json:"fee"`
	Status        string             `bson:"status" json:"status"`
}
