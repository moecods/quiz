package quiz

import (
	"encoding/json"
	"io"
	"moecods/quiz/utils"
	"net/http"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuizHandler struct {
	service QuizService
	repo    QuizRepository
}

func NewQuizHandler(service QuizService, repo QuizRepository) *QuizHandler {
	return &QuizHandler{repo: repo, service: service}
}

// GetQuiz godoc
//
//	@Summary		get list of quizzes
//	@Description	get list of quizzes
//	@Tags			quizzes
//	@Accept			json
//	@Produce		json
//	@Router			/quizzes [get]
func (h *QuizHandler) GetQuizzesHandler(w http.ResponseWriter, r *http.Request) {
	quizzes, _ := h.repo.ListQuizzes()
	utils.RespondWithJSON(w, http.StatusOK, quizzes)
}

// AddQuizHandler godoc
//
//	@Summary		Add a quiz
//	@Description	Add an quiz
//	@Tags			quizzes
//	@Accept			json
//	@Produce		json
//	@param			quiz	body		Quiz	true	"Quiz object"
//	@Failure		400		{string}	string	"Invalid request body"
//	@Failure		500		{string}	string	"Failed to read request body"
//	@Router			/quizzes [post]
func (h *QuizHandler) AddQuizHandler(w http.ResponseWriter, r *http.Request) {
	var quiz Quiz

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &quiz); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	quizValidation := QuizValidation{quiz: &quiz}
	if err := quizValidation.QuizValidation(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quiz, err = h.service.AddQuiz(quiz)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	utils.RespondWithJSON(w, http.StatusCreated, quiz)
}

// UpdateQuizHandler godoc
//
//	@Summary		Update a quiz
//	@Description	Update an existing quiz by ID
//	@Tags			quizzes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Quiz ID"
//	@Param			quiz	body		Quiz	true	"Quiz object"
//	@Success		200		{object}	Quiz
//	@Failure		400		{string}	string	"Invalid quiz ID or request body"
//	@Failure		500		{string}	string	"Internal Server Error"
//	@Router			/quizzes/{id} [put]
func (h *QuizHandler) UpdateQuizHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid quiz ID", http.StatusBadRequest)
		return
	}

	var quiz Quiz
	if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !quiz.ID.IsZero() && quiz.ID != id {
		http.Error(w, "Quiz ID in the request body does not match the ID in the URL", http.StatusBadRequest)
		return
	}

	quiz.ID = id

	if err := h.repo.UpdateQuiz(id, &quiz); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, quiz)
}

// DeleteQuizHandler godoc
//
//	@Summary		delete a quiz
//	@Description	delete an quiz
//	@Tags			quizzes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Quiz ID"
//	@Failure		400	{string}	string	"Invalid request body"
//	@Router			/quizzes/{id} [delete]
func (h *QuizHandler) DeleteQuizHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid quiz ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteQuiz(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Quiz deleted successfully"})
}

// GetQuiz godoc
//
//	@Summary		Show an quiz
//	@Description	get string by ID
//	@Tags			quizzes
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Quiz ID"
//	@Router			/quizzes/{id} [get]
func (h *QuizHandler) GetQuizHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid quiz ID", http.StatusBadRequest)
		return
	}

	quiz, err := h.service.GetQuiz(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, quiz)
}
