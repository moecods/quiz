package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"	
	"go.mongodb.org/mongo-driver/mongo"
)

type QuizRepository struct {
	collection *mongo.Collection
}

func NewQuizRepository(collection *mongo.Collection) *QuizRepository {
	return &QuizRepository{collection: collection}
}

func (r *QuizRepository) GetQuizzes() ([]Quiz, error)  {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := GetQuizCollection()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	var quizzes []Quiz
	for cursor.Next(ctx) {
		var quiz Quiz
		err := cursor.Decode(&quiz)
		if err != nil {
			log.Fatal(err)
		}
		quizzes = append(quizzes, quiz)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return quizzes, nil
}