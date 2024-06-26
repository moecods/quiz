package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnswerRepository struct {
	collection *mongo.Collection
}

func NewAnswerRepository(collection *mongo.Collection) *AnswerRepository {
	return &AnswerRepository{collection: collection}
}

func (r *AnswerRepository) ListAnswers() ([]Answer, error)  {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	var answers []Answer
	for cursor.Next(ctx) {
		var answer Answer
		err := cursor.Decode(&answer)
		if err != nil {
			log.Fatal(err)
		}
		answers = append(answers, answer)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return answers, nil
}

func (r *AnswerRepository) AddAnswer(answer *Answer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	result, err := r.collection.InsertOne(ctx, answer)
	answer.ID = result.InsertedID.(primitive.ObjectID)

	log.Println(result)

	return err
}

func (r *AnswerRepository) UpdateAnswer(id primitive.ObjectID, answer *Answer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": answer,
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *AnswerRepository) DeleteAnswer(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *AnswerRepository) GetAnswer(id primitive.ObjectID) (*Answer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var answer Answer
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&answer)
	return &answer, err
}