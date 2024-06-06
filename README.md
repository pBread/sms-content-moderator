# SMS Content Moderator

## Overview

The SMS Content Moderator is a service designed to help businesses monitor and control the content of SMS messages. This service allows administrators to configure content moderation policies to evaluate messages for potentially inappropriate or restricted content.

_Note: This app requires configuration for individual use cases (see [Configuring Content Rules](#configuring-content-rules)). This app does not gaurantee messages are compliant with SMS guidelines; rather, it serves as a supportive tool in compliance efforts._

## Quickstart

### Prerequisites

- [Install Go on your machine](https://go.dev/doc/install)
- [Create an OpenAI Platform Account](https://platform.openai.com/signup)

#### Clone Repo

```bash
git clone https://github.com/yourgithub/sms-content-moderator.git
cd sms-content-moderator
```

#### Set Environment Variables

```bash
cp .env.example .env
```

Open the .env file in a text editor and set the following variables:

- OPENAI_API_KEY: Your OpenAI API key.
- OPENAI_MODEL: The model identifier you intend to use (e.g., "gpt-4").

#### Run Application

```bash
go run ./cmd/server
```

#### Test the API

```bash
curl -X POST \
     -H "Content+Type: application/json" \
     -d '{"Message": "This message contains a maybe bad word."}' \
     http://localhost:8080/evaluate-message
```

## Configuring Content Rules

# DEPRECATED

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

CSV must include a header.

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
