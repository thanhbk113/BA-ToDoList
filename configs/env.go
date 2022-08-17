package configs

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

//func EnvMongoURI get MONGOURI from .env file
func EnvMongoURI() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println(os.Getenv("MONGOURI"))

	return os.Getenv("MONGOURI")
}

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI())) // create new client with mongo uri

	if err != nil {
		log.Fatal("err connect mongo db uri:", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //get context and cancel context with 10 second

	err = client.Connect(ctx) // checck conntect success or not

	if err != nil {
		log.Fatal("err connect MongoDb:", err)
	}

	//Ping to Database make sure connect database success

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal("Err pins to MongoDb")
	}

	fmt.Println("Connected MongoDb")

	return client
}

//Get client from func ConnectDB
var DB *mongo.Client = ConnectDB()

//getting database collections

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("golangApi").Collection(collectionName)

	return collection
}
