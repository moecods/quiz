package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"moecods/quiz/utils"
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
	utils.RespondWithJSON(w, http.StatusOK, answers)
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
	utils.RespondWithJSON(w, http.StatusCreated, answerRequest)
}
