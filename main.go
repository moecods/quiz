package main

import (
	"context"
	"log"
	"net/http"
	"moecods/quiz/quiz"
	"moecods/quiz/participant"
	"github.com/gorilla/mux"

)

func main() {
	client := ConnectToDB()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	quizCollection := GetQuizCollection()
	quizRepo := quiz.NewQuizRepository(quizCollection)
	quizService := *quiz.NewQuizService(quizRepo)
	quizHandler := quiz.NewQuizHandler(quizService, quizRepo)

	r := mux.NewRouter()
	r.HandleFunc("/quizzes", recoverHandler(quizHandler.GetQuizzesHandler)).Methods(http.MethodGet)
	r.HandleFunc("/quizzes", recoverHandler(quizHandler.AddQuizHandler)).Methods(http.MethodPost)
	r.HandleFunc("/quizzes/{id}", recoverHandler(quizHandler.UpdateQuizHandler)).Methods(http.MethodPut)
	r.HandleFunc("/quizzes/{id}", recoverHandler(quizHandler.DeleteQuizHandler)).Methods(http.MethodDelete)
	r.HandleFunc("/quizzes/{id}", recoverHandler(quizHandler.GetQuizHandler)).Methods(http.MethodGet)

	participantCollection := GetParticipantCollection()
	participantRepo := participant.NewParticipantRepository(participantCollection)
	participantHandler := participant.NewParticipantHandler(*participantRepo)

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
