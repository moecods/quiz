package participant

import (
	"encoding/json"
	"io"
	"log"
	"moecods/quiz/utils"
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

type ParticipantAnswer struct {
	ParticipantId primitive.ObjectID `json:"participant_id"`
	Answers       []Answer           `json:"answers"`
}

type SaveParticipantsAnswersRequest struct {
	ParticipantAnswers []ParticipantAnswer `json:"participant_answers" bson:"participant_answers"`
}

func NewParticipantHandler(repo ParticipantRepository) *ParticipantHandler {
	return &ParticipantHandler{ParticipantRepo: repo}
}

// RegisterParticipantsHandler godoc
//
//	@Summary		Register participants
//	@Description	Register participants
//	@Tags			participants
//	@Accept			json
//	@Produce		json
//	@Param			RegisterRequest	body		RegisterRequest	true	"Register Request object"
//	@Failure		400				{string}	string			"Invalid request body"
//	@Failure		500				{string}	string			"Failed to read request body"
//	@Router			/participants/register [post]
func (h *ParticipantHandler) RegisterParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var RegisterRequest RegisterRequest

	if err := json.Unmarshal(body, &RegisterRequest); err != nil {
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
		http.Error(w, "Failed to store", http.StatusInternalServerError)
	}
	utils.RespondWithJSON(w, http.StatusOK, participants)
}

// RegisterParticipantsHandler godoc
//
//	@Summary		Save participants answers
//	@Description	Save participants answers
//	@Tags			participants
//	@Accept			json
//	@Produce		json
//	@Param			SaveParticipantsAnswersRequest	body		SaveParticipantsAnswersRequest	true	"Save participants answers Request object"
//	@Failure		404								{string}	string							"Participant Not Found"
//	@Failure		400								{string}	string							"Invalid request body"
//	@Failure		500								{string}	string							"Failed to read request body"
//	@Router			/participants/answers [post]
func (h *ParticipantHandler) SaveParticipantsAnswersHandler(w http.ResponseWriter, r *http.Request) {
	var participantsAnswersRequest SaveParticipantsAnswersRequest

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &participantsAnswersRequest); err != nil {
		log.Printf("Failed to unmarshal request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, particpantAnswersRequest := range participantsAnswersRequest.ParticipantAnswers {
		participant, err := h.ParticipantRepo.GetParticipant(particpantAnswersRequest.ParticipantId)
		if err != nil {
			http.Error(w, "Participant Not Found", 404)
		}

		for i, answerRequest := range particpantAnswersRequest.Answers {
			existingAnswer, found := findAnswerByQuestionID(participant.Answers, answerRequest.QuestionID)
			if found {
				existingAnswer.AnswerText = answerRequest.AnswerText
				existingAnswer.SelectedOption = answerRequest.SelectedOption
				existingAnswer.IsCorrect = answerRequest.IsCorrect
				existingAnswer.AnsweredAt = answerRequest.AnsweredAt
				participant.Answers[i] = existingAnswer
			} else {
				newAnswer := Answer{
					QuestionID:     answerRequest.QuestionID,
					AnswerText:     answerRequest.AnswerText,
					SelectedOption: answerRequest.SelectedOption,
					IsCorrect:      answerRequest.IsCorrect,
					AnsweredAt:     answerRequest.AnsweredAt,
				}
				participant.Answers = append(participant.Answers, newAnswer)
			}
		}

		err = h.ParticipantRepo.UpdateParticipant(participant.ID, participant)

		if err != nil {
			http.Error(w, "Failed to update participant", http.StatusInternalServerError)
		}
	}

	utils.RespondWithJSON(w, http.StatusCreated, struct{}{})
}

func findAnswerByQuestionID(answers []Answer, QuestionID primitive.ObjectID) (Answer, bool) {
	for _, answer := range answers {
		if answer.QuestionID == QuestionID {
			return answer, true
		}
	}
	return Answer{}, false
}
