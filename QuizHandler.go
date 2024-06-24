package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)


func GetQuizzes(w http.ResponseWriter, r *http.Request) {
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

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(quizzes); err != nil {
        log.Printf("Failed to encode response: %v", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}