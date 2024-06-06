# SMS Content Moderator

The SMS Content Moderator is a service designed to help businesses monitor and control the content of SMS messages. This service allows administrators to configure content moderation policies to evaluate messages for potentially inappropriate or restricted content.

_Note: This app requires configuration for individual use cases (see [Configuring Content Rules](#configuring-content-rules)). This app does not gaurantee messages are compliant with SMS guidelines; rather, it serves as a supportive tool in compliance efforts._

### How the App Works

The SMS Content Moderator operates through a straightforward yet effective process designed to ensure SMS compliance with relevant guidelines:

1. **Configuration by Administrators**: Administrators must first configure the system by setting up a blacklist.csv and corresponding policy documents for each policy. The blacklist entries specify patterns to match (either as direct strings or regex) and categorize them by severity (Tier 0 for automatic rejection, Tier 1 for further review).

2. **Blacklist and Policy Matching**: When a message is received via API, the app scans the content against the blacklist entries. If a match is found, the response depends on the tier:

- **Tier 0**: Messages matching these entries are immediately flagged as violations, and the API response includes the specific policies breached. Such messages are recommended for rejection.
- **Tier 1**: These entries trigger a deeper examination. The content is further analyzed using policy documents related to the matched entries to determine the context and intent.

3. **Contextual Analysis with LLM**: For Tier 1 matches, the app compiles relevant policies into a prompt and consults an LLM (like OpenAI) to assess if the message content indeed violates the intended policies. This step ensures that messages are not wrongly flagged based on out-of-context words or phrases.

This system allows businesses to customize their moderation tools extensively, ensuring that SMS content aligns with both internal standards and regulatory requirements.

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

Configuration of the SMS Content Moderator involves setting up [blacklist](config/blacklist.csv) entries and [policy documents](config/policies/) to define what content is checked and how it is evaluated.

_Important: The provided blacklist and policy documents serve as examples and must be customized to the intended use-case._

### Blacklist Configuration

- **Blacklist File**: The blacklist is defined in a CSV located at [config/blacklist.csv](config/blacklist.csv).
- **CSV Format**: The CSV file <u>must include a header row</u> with the columns in this order: `Content`, `Content Type`, `Policy`, `Tier`. Each row represents one blacklist entry.

  - **Content**: The text or regex pattern to match against. This app uses the [Go regex syntax](https://pkg.go.dev/regexp/syntax).
  - **Content Type**: `regex` for regular expressions, `string` for direct string matches
  - **Policy**: The name of the policy, which must correspond to a policy markdown file in the [config/policies](config/policies) directory.
  - **Tier**: `0` for words that trigger auto-rejection, `1` for words that indicate a possible violation but require further LLM evaluation.

### Policy Documents

- **Location and Naming**: Each policy referred to in the `Policy` column of the blacklist must have a corresponding markdown document in the [config/policies](config/policies) directory. For example, if a blacklist entry has the policy `profanity` then there must be a document describing that policy located here: [config/policies/profanity.md](config/policies/profanity.md)

- **Customization**: You are encouraged to review and modify the provided policy documents to fit your use-case. You can also create new policies by adding corresponding entries to the blacklist CSV and creating new policy markdown files.

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
