package main

import (
	"context"
	"log"
	_ "moecods/quiz/docs"
	"moecods/quiz/participant"
	"moecods/quiz/quiz"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.email	moecods.dev@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8020
// @BasePath	/v1
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
	registerSwagger(r)

	r.HandleFunc("/v1/quizzes", recoverHandler(quizHandler.GetQuizzesHandler)).Methods(http.MethodGet)
	r.HandleFunc("/v1/quizzes", recoverHandler(quizHandler.AddQuizHandler)).Methods(http.MethodPost)
	r.HandleFunc("/v1/quizzes/{id}", recoverHandler(quizHandler.UpdateQuizHandler)).Methods(http.MethodPut)
	r.HandleFunc("/v1/quizzes/{id}", recoverHandler(quizHandler.DeleteQuizHandler)).Methods(http.MethodDelete)
	r.HandleFunc("/v1/quizzes/{id}", recoverHandler(quizHandler.GetQuizHandler)).Methods(http.MethodGet)

	participantCollection := GetParticipantCollection()
	participantRepo := participant.NewParticipantRepository(participantCollection)
	participantHandler := participant.NewParticipantHandler(*participantRepo)

	r.HandleFunc("/v1/participants/answers", recoverHandler(participantHandler.SaveParticipantsAnswersHandler)).Methods(http.MethodPost)
	r.HandleFunc("/v1/participants/register", recoverHandler(participantHandler.RegisterParticipantsHandler)).Methods(http.MethodPost)
	r.HandleFunc("/v1/quizzes/{id}/participants", recoverHandler(participantHandler.GetParticipantsByQuizIDHandler)).Methods(http.MethodGet)

	http.Handle("/", r)
	err := http.ListenAndServe(":8020", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
	log.Println("Starting server on :8020")
}

func registerSwagger(r *mux.Router) {
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8020/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
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
