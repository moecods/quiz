package main

import (
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnswerHandler struct {
    AnswerRepo AnswerRepository
}

func  NewAnswerHandler(repo AnswerRepository) *AnswerHandler {
    return &AnswerHandler{AnswerRepo: repo}
}

type AnswerRequest struct {
    Answers []Answer `json:"answers"`
}

func (h *AnswerHandler) GetAnswersHandler(w http.ResponseWriter, r *http.Request) {
	answers, _ := h.AnswerRepo.ListAnswers()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(answers); err != nil {
        log.Printf("Failed to encode response: %v", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}


func (h *AnswerHandler) AddAnswersHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
	idStr := vars["id"]
	
	quizId, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid quiz ID", http.StatusBadRequest)
		return
	}

    var answerRequest AnswerRequest

    body, err := ioutil.ReadAll(r.Body)

    if err != nil {
        log.Printf("Failed to read request body: %v", err)
        http.Error(w, "Failed to read request body", http.StatusInternalServerError)
        return
    }

    if err := json.Unmarshal(body, &answerRequest); err != nil {
        log.Printf("Failed to unmarshal request body: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    for i := range answerRequest.Answers {
        answerRequest.Answers[i].QuizID = quizId
        err = h.AnswerRepo.AddAnswer(&answerRequest.Answers[i])
        if err != nil {
            http.Error(w, err.Error() , http.StatusBadRequest)
        }
    }

	
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated) // 201 Created
    if err := json.NewEncoder(w).Encode(answerRequest.Answers); err != nil {
        log.Printf("Failed to encode response: %v", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}