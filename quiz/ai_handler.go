package quiz

import (
	"context"
	"encoding/json"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"moecods/quiz/utils"
	"net/http"
	"os"
)

type AiHandler struct {
}

func NewAiHandler() *AiHandler {
	return &AiHandler{}
}

type AIRequest struct {
	Message string `json:"message"`
}

// GenerateAIResponse godoc
//
//	@Summary		Generate AI answer for a given question
//	@Description	Accepts a JSON payload with a user message and returns AI-generated response.
//	@Tags			AI
//	@Accept			json
//	@Produce		json
//	@Param			request	body		AIRequest	true	"User question message"
//	@Success		200		{object}	map[string]string	"AI generated answer"
//	@Failure		400		{object}	map[string]string	"Invalid input or AI error"
//	@Router			/ai/answer [post]
func (h *AiHandler) GenerateAIResponse(w http.ResponseWriter, r *http.Request) {
	const model = "openai/gpt-4o-mini"

	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
		option.WithBaseURL(os.Getenv("OPENAI_BASE_URL")),
	)

	var req AIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(req.Message),
		},
		Model: model,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, chatCompletion)
}
