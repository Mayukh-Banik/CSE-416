package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 구조체는 사용자 정보를 정의합니다.
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // MongoDB에서 사용하는 ObjectID
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	PublicKey string             `bson:"public_key" json:"public_key"`
	Balance   float64            `bson:"balance" json:"balance"`
	Reputation int               `bson:"reputation" json:"reputation"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
