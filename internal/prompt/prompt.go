package prompt

import (
	"fmt"
	"os"
	"strings"
)

// BuildPrompt constructs the prompt for content evaluation based on the base prompt and policies.
func BuildPrompt(content string, policies []string) (string, error) {
	// Read the base prompt from config/prompt.txt
	basePrompt, err := os.ReadFile("config/prompt.txt")
	if err != nil {
		return "", fmt.Errorf("failed to read base prompt: %v", err)
	}

	// Initialize variables to be injected into the prompt
	policyNames := []string{}
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

		policyNames = append(policyNames, policyName)

		policyFilePath := fmt.Sprintf("config/policies/%s.txt", policyName)
		if _, err := os.Stat(policyFilePath); err == nil {
			policyContent, err := os.ReadFile(policyFilePath)
			if err == nil {
				policyNotes += fmt.Sprintf("\n\n%s:\n%s", policyName, string(policyContent))
			}
		}
	}

	// Inject the content and policy variables into the base prompt
	policiesList := strings.Join(policyNames, ", ")
	prompt := string(basePrompt)
	prompt = strings.Replace(prompt, "{{content}}", content, 1)
	prompt = strings.Replace(prompt, "{{policies}}", policiesList, 1)
	prompt = strings.Replace(prompt, "{{policyNotes}}", policyNotes, 1)

	return prompt, nil
}
