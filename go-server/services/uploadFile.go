package services 

import(
	"context"
	"log"
	"go-server/models"
	"go-server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const fileCollection = "files"

func StoreFileMetaData(file models.File) error 
{
	collection := utils.GetCollection("squidcoinDB", fileCollection)
    _, err = collection.InsertOne(context.TODO(), file)

	if err != nil{
		log.Printf("failed to insert metadata into MongoDB:%v",err)
		return err 
	}
	return nil 
}