package ai

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type AIClient interface {
	GetChatCompletion(ctx context.Context, prompt string) (string, error)
}

type LiaraAIClient struct {
	client *openai.Client
	model  string
}

func NewLiaraAIClient(apiKey, baseURL, model string) AIClient {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	return &LiaraAIClient{
		client: client,
		model:  model,
	}
}

func (c *LiaraAIClient) GetChatCompletion(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
