package client

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type AIClient struct {
	client *openai.Client
}

func NewAIClient(token string) *AIClient {
	return &AIClient{
		client: openai.NewClient(token),
	}
}

func (ai *AIClient) QuestionAI(text string) (string, error) {
	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo16K0613,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("error ai.client.CreateChatCompletion(): %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
