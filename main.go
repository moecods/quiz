package main

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	ID primitive.ObjectID
	Title string
	Description string
	Questions []Question
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Question struct {
	title string
	correct_answer string
	user_answer string
}

func main() {
	quiz := Quiz{
		Title: "english A1",
		Description: "General Knowledge English",
		CreatedAt:  time.Now(),
		UpdatedAt: time.Now(),
	} 

	fmt.Println(quiz)
}