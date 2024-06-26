package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type          string             `bson:"type" json:"type"` // Type of question: "descriptive" or "test"
	Text          string             `bson:"text" json:"text"`
	Options       []string           `bson:"options" json:"options"`
	CorrectOption int                `bson:"correct_option" json:"correct_option"`
}

type Answer struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	QuizID         primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	QuestionID     primitive.ObjectID `bson:"question_id" json:"question_id"`
	Type           string             `bson:"type" json:"type"`
	AnswerText     string             `bson:"answer_text" json:"answer_text"`
	SelectedOption int				  `bson:"selection_option" json:"selection_option"`
	IsCorrect      bool               `bson:"is_correct" json:"is_correct"`
	AnsweredAt     time.Time          `bson:"answered_at" json:"answered_at"`
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

	r := mux.NewRouter()
	r.HandleFunc("/quizzes", recoverHandler(quizHandler.GetQuizzesHandler)).Methods(http.MethodGet)
	r.HandleFunc("/quizzes", recoverHandler(quizHandler.AddQuizHandler)).Methods(http.MethodPost)
	r.HandleFunc("/quizzes/{id}", recoverHandler(quizHandler.UpdateQuizHandler)).Methods(http.MethodPut)
	r.HandleFunc("/quizzes/{id}", recoverHandler(quizHandler.DeleteQuizHandler)).Methods(http.MethodDelete)
	r.HandleFunc("/quizzes/{id}", recoverHandler(quizHandler.GetQuizHandler)).Methods(http.MethodGet)

	answerCollection := GetAnswerCollection()
	answerRepo := NewAnswerRepository(answerCollection)
	answerHandler := NewAnswerHandler(*answerRepo)

	r.HandleFunc("/quizzes/{id}/answers", recoverHandler(answerHandler.AddAnswersHandler)).Methods(http.MethodPost)

	http.Handle("/", r)
	err := http.ListenAndServe(":8020", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
	log.Println("Starting server on :8020")
}

func recoverHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Recovered from panic: %v", rec)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		h(w, r)
	}
}
