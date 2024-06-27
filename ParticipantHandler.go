package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParticipantHandler struct {
	ParticipantRepo ParticipantRepository
}

type RegisterRequest struct {
	QuizId               primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	NumberOfParticipants int                `bson:"number_of_participants" json:"number_of_participants"`
}

func NewParticipantHandler(repo ParticipantRepository) *ParticipantHandler {
	return &ParticipantHandler{ParticipantRepo: repo}
}

func (h *ParticipantHandler) RegisterParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var RegisterRequest RegisterRequest

	if err := json.Unmarshal(body, &RegisterRequest); err != nil {
		log.Printf("Failed to unmarshal request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var participants []Participant
	for i := 0; i < RegisterRequest.NumberOfParticipants; i++ {
		participants = append(participants, Participant{
			ID:     primitive.NewObjectID(),
			QuizID: RegisterRequest.QuizId,
			Status: "not_started",
		})
	}

	err = h.ParticipantRepo.AddManyParticipants(participants)

	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to store", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(participants); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
