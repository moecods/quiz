package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnswerHandler struct {
	AnswerRepo AnswerRepository
}

func NewAnswerHandler(repo AnswerRepository) *AnswerHandler {
	return &AnswerHandler{AnswerRepo: repo}
}

type ParticipantAnswer struct {
	ParticipantId primitive.ObjectID `json:"participant_id"`
	Answers       []Answer           `json:"answers"`
}

type SaveAnswerRequest struct {
	ParticipantAnswers []ParticipantAnswer `json:"participant_answers" bson:"participant_answers"`
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
	var answerRequest SaveAnswerRequest

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// log.Printf("Request body: %s", string(body))
	
	if err := json.Unmarshal(body, &answerRequest); err != nil {
		log.Printf("Failed to unmarshal request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created

	if err := json.NewEncoder(w).Encode(answerRequest); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
