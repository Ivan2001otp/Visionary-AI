package database

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client

var lockVariable = &sync.Mutex{}

const mongoUrl string = "mongodb://localhost:27017"

func connectToDB() error {
	clientOptions := options.Client().ApplyURI(mongoUrl)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		fmt.Println("Something went wrong while connecting to mongo db in DbConfig.go")
		log.Fatal(err)
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	db = client
	fmt.Println("Successfully connected to database")
	return nil
}

func MongoDbInstanceProvider() (*mongo.Client, error) {
	if db == nil {
		//locks only for single go routine to create instance..during concurrent operations.
		lockVariable.Lock()

		defer lockVariable.Unlock()
		if db == nil {
			fmt.Println("Creating singleton inst")
			err := connectToDB()

			if err != nil {
				return nil, err
			}
		} else {
			fmt.Println("Singleton instance already provided..")
		}
	} else {
		fmt.Println("Singleton instance already provided..")
	}

	return db, nil
}
