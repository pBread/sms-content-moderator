package llm

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// BuildPrompt constructs the prompt for content evaluation based on the base prompt and policies.
func BuildPrompt(content string, policies []string) (string, error) {
	// Read the base prompt from config/prompt.md
	basePrompt, err := os.ReadFile("config/prompt.md")
	if err != nil {
		return "", fmt.Errorf("failed to read base prompt: %v", err)
	}

	// Initialize variables to be injected into the prompt
	policyNotes := ""

	// Process each policy and gather the policy notes
	for _, policy := range policies {
		var policyName string
		policyParts := strings.SplitN(policy, "-", 2)
		if len(policyParts) == 2 {
			policyName = policyParts[1]
		} else {
			policyName = policyParts[0]
		}

		policyName = strings.TrimSpace(policyName)

		policyFilePath := fmt.Sprintf("config/policies/%s.md", policyName)
		if _, err := os.Stat(policyFilePath); err == nil {
			policyContent, err := os.ReadFile(policyFilePath)
			if err == nil {
				policyNotes += fmt.Sprintf("\n\n===POLICY_ID:'%s'===\n%s", policyName, string(policyContent))
			}
		}
	}

	// Inject the content and policy variables into the base prompt
	prompt := string(basePrompt)
	prompt = strings.Replace(prompt, "{{content}}", content, 1)
	prompt = strings.Replace(prompt, "{{policies}}", policyNotes, 1)

	return prompt, nil
}

func EvalPolicyViolation(content string) (string, error) {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	openaiModel := os.Getenv("OPENAI_MODEL")

	client := openai.NewClient(openaiKey)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openaiModel,
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