package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Questions   []Question         `bson:"questions" json:"questions"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Question struct {
	ID            primitive.ObjectID
	Type          string
	Text          string
	Options       []string
	CorrectOption int
}

type Answer struct {
	ID             primitive.ObjectID
	QuizID         primitive.ObjectID
	QuestionID     primitive.ObjectID
	Type           string
	AnswerText     string
	SelectedOption int
	IsCorrect      bool
	AnsweredAt     time.Time
}

func main() {
	quiz := Quiz{
		ID:          primitive.NewObjectID(),
		Title:       "english A1",
		Description: "General Knowledge English",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	questions := []Question{
		{
			ID:   primitive.NewObjectID(),
			Type: "descriptive",
			Text: "Describe the process of photosynthesis.",
		},
		{
			ID:            primitive.NewObjectID(),
			Type:          "test",
			Text:          "What is the capital of France?",
			Options:       []string{"Paris", "London", "Berlin", "Madrid"},
			CorrectOption: 0,
		},
	}

	quiz.Questions = questions

	var answers []Answer

	answerText := "Photosynthesis is the process by which green plants and some other organisms use sunlight to synthesize foods with the help of chlorophyll."
	answers = append(answers, Answer{
		ID:         primitive.NewObjectID(),
		QuizID:     quiz.ID,
		QuestionID: quiz.Questions[0].ID,
		AnswerText: answerText,
		AnsweredAt: time.Now(),
	})

	answers = append(answers, Answer{
		ID:             primitive.NewObjectID(),
		QuizID:         quiz.ID,
		QuestionID:     quiz.Questions[1].ID,
		SelectedOption: 1,
		IsCorrect:      false,
		AnsweredAt:     time.Now(),
	})

	client := ConnectToDB()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	quizCollection := GetQuizCollection()
	quizRepo := NewQuizRepository(quizCollection)
	quizHandler := NewQuizHandler(*quizRepo)

	http.HandleFunc("GET /quizzes/", quizHandler.GetQuizzesHandler)
	log.Println("Starting server on :8020")
	err := http.ListenAndServe(":8020", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
