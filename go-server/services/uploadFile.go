package services

import (
	"context"
	"go-server/models"
	"go-server/utils"
	"log"
)

const fileCollection = "files"

func StoreFileMetaData(file models.File) error {
	collection := utils.GetCollection("squidcoinDB", fileCollection)
	_, err := collection.InsertOne(context.TODO(), file)

	if err != nil {
		log.Printf("failed to insert metadata into MongoDB:%v", err)
		return err
	}
	return nil
}
