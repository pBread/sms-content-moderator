# SMS Content Moderator

## Docs to Add

### Blacklist CSV Definition

A CSV with this definition, in that specific order.

```golang
type CSVBlacklistEntry struct {
	Content     string // the text or regex pattern to match against
	ContentType string // "regex" for regular expressions, "string" for direct string matches
	Policy      string // a descriptor of the policy the entry enforces
	Tier        int    // importance level: 0 for auto-rejection, 1 for LLM evaluation
}
```

### Prompt Definition

Base prompt is defined in config/prompt.txt.

Prompts are appended with additional information based on which policies are violated (see ### Policy Definition)

Describe how to update prompts. Explain how the prompt is assembled and how it is sent to the LLM.

### Policy Definition

Every policy can have an specific prompt that may including a .txt file in config/policies/[policy].txt

Create a section on how to update policies

### Deployment Examples

The directory cmd/ has multiple examples of various deployments.

- cmd/server is a simple server. Note: It doesn't have authentication
- cmd/conversations shows a simple server example deployed for Twilio Conversations.

### Twilio Conversations Example

- Twilio is authenticated through the X-Twilio-Signature, see https://www.twilio.com/docs/usage/webhooks/webhooks-security#validating-signatures-from-twilio

- You need to define the env variables CONVERSATIONS_PRE_EVENT_WEBHOOK_URL & TWILIO_AUTH_TOKEN

### Explain the prompt can be offloaded to the LLM provider, most likely
