package prompt

import (
	"fmt"
	"os"
	"strings"
)

func BuildPrompt(content string, policies []string) (string, error) {
	// Read the base prompt from config/prompt.txt
	basePrompt, err := os.ReadFile("config/prompt.txt")
	if err != nil {
		return "", fmt.Errorf("failed to read base prompt: %v", err)
	}

	// Initialize the prompt with the base content
	prompt := string(basePrompt)
	prompt = strings.Replace(prompt, "{{content}}", content, 1)

	// Add each policy-specific file content if it exists
	for _, policy := range policies {
		policyParts := strings.SplitN(policy, "-", 2)
		if len(policyParts) != 2 {
			continue
		}
		policyFilePath := fmt.Sprintf("config/policies/%s.txt", policyParts[1])
		if _, err := os.Stat(policyFilePath); err == nil {
			policyContent, err := os.ReadFile(policyFilePath)
			if err == nil {
				prompt += fmt.Sprintf("\n\n%s:\n%s", policyParts[1], string(policyContent))
			}
		}
	}

	return prompt, nil
}
