package prompts

var basePrompt = `
Determine whether the message violates the content policy.

Message:

Policy:

Additional Context:



Respond with the following JSON format:
{
	isViolation: boolean
	reason: string
}



`

var PromptMap = map[string]string{
	"abuse":        ``,
	"alcohol":      ``,
	"cannabis":     ``,
	"firearms":     ``,
	"gambling":     ``,
	"misc":         ``,
	"prescription": ``,
	"profanity":    ``,
	"tobacco":      ``,
}
