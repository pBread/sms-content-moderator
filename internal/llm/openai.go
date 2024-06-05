package llm

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func EvalPolicyViolation(content string) (string, error) {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(openaiKey)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	log.Println(resp.Choices[0].Message.Content)

	return resp.Choices[0].Message.Content, nil

}
