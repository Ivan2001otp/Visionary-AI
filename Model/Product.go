package model

import (
	"context"
	"fmt"
	"log"

	constant "github.com/Ivan2001otp/Visionary-AI/Constant"
	config "github.com/Ivan2001otp/Visionary-AI/Service/Database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductName      string             `json:"productname" bson:"productname"`
	ProductImg       string             `json:"productimg" bson:"productimg"`
	ProductDetailUrl string             `json:"productdetailurl" bson:"productdetailurl"`
	ProductRating    string             `json:"productrating" bson:"productrating"`
	GlobalRating     string             `json:"globalrating" bson:"globalrating"`
	ProductPrice     string             `json:"productprice" bson:"productprice"`
	ProductRetailers string             `json:"productretailers" bson:"productretailers"`
	CategoryType     string             `json:"categorytype" bson:"categorytype"`
}

/*
this is an interface that has abstract methods
to perform CRUD operations.
*/
type iDbMethod interface {
	SaveToMongo(product Product) (bool, error)
	DeleteByProductId(id string) (bool, error)
	UpdateByProductId(id string) (bool, error)
	fetchAllProduct() ([]Product, error)
	DeleteAllFromMongo() (bool, error)
}

// implementing the interface...
func DeleteAllFromMongo() (int, error) {
	//this deletes all the records from mongo..
	fmt.Println("Initiating delete all function")
	mongoClient, err := config.MongoDbInstanceProvider()

	if err != nil {
		log.Fatal("Error while deleting all product->", err)
	}

	productCollection := mongoClient.Database(constant.DATABASE_NAME).Collection(constant.COLLECTION)

	status, err := productCollection.DeleteMany(context.TODO(), bson.M{})

	if err != nil {
		return -1, err
	}

	return int(status.DeletedCount), nil

}

func (p Product) SaveToMongo() (bool, error) {
	fmt.Println("The product saved is ", p)

	mongoClient, err := config.MongoDbInstanceProvider()
	if err != nil {
		log.Fatal("error while saving..", err)
	}

	productCollection := mongoClient.Database(constant.DATABASE_NAME).Collection(constant.COLLECTION)
	status, err := productCollection.InsertOne(context.TODO(), p)

	if err != nil {
		return false, err
	}
	p.Id = status.InsertedID.(primitive.ObjectID)

	fmt.Println("Saved")

	return true, nil
}
