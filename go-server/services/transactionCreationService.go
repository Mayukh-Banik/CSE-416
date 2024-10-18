// services/transactionCreationService.go
package services

import (
	"errors"
	"time"
	"go-server/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewTransaction 함수는 트랜잭션을 생성하고 sender와 receiver가 다른지 확인합니다.
func NewTransaction(sender, receiver primitive.ObjectID, amount float64) (*models.Transaction, error) {
	if sender == receiver {
		return nil, errors.New("sender and receiver cannot be the same user")
	}

	return &models.Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Timestamp: time.Now(),
		Status:    "pending", // 기본값으로 'pending' 상태
	}, nil
}
