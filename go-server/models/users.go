package models

import "time"

type User struct {
    UserID      string    `bson:"user_id" json:"user_id"`
    PublicKey   string    `bson:"public_key" json:"public_key"`
    CreatedDate time.Time `bson:"created_date" json:"created_date"`
}
