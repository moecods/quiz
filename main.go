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

type Participant struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	QuizID     primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	Status     string             `bson:"status" json:"status"` // not_started, started, finished
	Answers    []Answer           `bson:"answers" json:"answers"`
	StartAt    time.Time
	FinishedAt time.Time
}

type Answer struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	QuestionID     primitive.ObjectID `bson:"question_id" json:"question_id"`
	AnswerText     string             `bson:"answer_text" json:"answer_text"`
	SelectedOption int                `bson:"selection_option" json:"selection_option"`
	IsCorrect      bool               `bson:"is_correct" json:"is_correct"`
	AnsweredAt     time.Time          `bson:"answered_at" json:"answered_at"`
}

func main() {
	client := ConnectToDB()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	quizCollection := GetQuizCollection()
	quizRepo := NewQuizRepository(quizCollection)
	quizService := *NewQuizService(quizRepo)
	quizHandler := NewQuizHandler(quizService, quizRepo)

	r := mux.NewRouter()
	r.HandleFunc("/quizzes", recoverHandler(quizHandler.GetQuizzesHandler)).Methods(http.MethodGet)
	r.HandleFunc("/quizzes", recoverHandler(quizHandler.AddQuizHandler)).Methods(http.MethodPost)
	r.HandleFunc("/quizzes/{id}", recoverHandler(quizHandler.UpdateQuizHandler)).Methods(http.MethodPut)
	r.HandleFunc("/quizzes/{id}", recoverHandler(quizHandler.DeleteQuizHandler)).Methods(http.MethodDelete)
	r.HandleFunc("/quizzes/{id}", recoverHandler(quizHandler.GetQuizHandler)).Methods(http.MethodGet)

	participantCollection := GetParticipantCollection()
	participantRepo := NewParticipantRepository(participantCollection)
	participantHandler := NewParticipantHandler(*participantRepo)

	r.HandleFunc("/participants/answers", recoverHandler(participantHandler.SaveParticipantsAnswersHandler)).Methods(http.MethodPost)
	r.HandleFunc("/participants/register", recoverHandler(participantHandler.RegisterParticipantsHandler)).Methods(http.MethodPost)

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
