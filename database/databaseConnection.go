package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	// Loading the environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error in loading .env file")
	}
	// Getting mongo url
	MongoDb := os.Getenv("MONGODB_ATLAS_URL")
	// Creating new client with mongo url
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))

	if err != nil {
		log.Fatal(err)
	}
	// Creating context for db
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	// Connecting db with context
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongo db")

	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = (*mongo.Collection)(client.Database("cluster0").Collection(collectionName))
	return collection
}
