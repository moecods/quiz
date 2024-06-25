package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type QuizHandler struct {
    QuizRepo QuizRepository
}

func  NewQuizHandler(repo QuizRepository) *QuizHandler {
    return &QuizHandler{QuizRepo: repo}
}


func (h *QuizHandler) GetQuizzesHandler(w http.ResponseWriter, r *http.Request) {
	quizzes, _ := h.QuizRepo.GetQuizzes()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(quizzes); err != nil {
        log.Printf("Failed to encode response: %v", err)
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}