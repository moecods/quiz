package main

import (
	"context"
	"log"
	"fmt"
	"os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
    mongoClient *mongo.Client
    dbName      = "quiz_db"
)

func ConnectToDB()  *mongo.Client {
	uri := "mongodb://localhost:27017"

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

	err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
		os.Exit(1)
    }

	fmt.Println("Connected to MongoDB!")
	mongoClient = client
	return client
}

func GetQuizCollection() *mongo.Collection {
    if mongoClient == nil {
        ConnectToDB()
    }

    return mongoClient.Database(dbName).Collection("quizzes")
}