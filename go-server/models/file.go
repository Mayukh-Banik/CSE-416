package models

import (
    "time"
)

type File struct {
    FileName   string    `bson:"file_name" json:"file_name"`
    Hash       string    `bson:"hash" json:"hash"`
    Reputation int       `bson:"reputation" json:"reputation"`
    FileSize   int64     `bson:"file_size" json:"file_size"` // byte
    CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}
