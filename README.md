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

- `OPENAI_API_KEY`: Your OpenAI API key.
- `OPENAI_MODEL`: The model identifier you intend to use (e.g., "gpt-4").

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
- **CSV Format**: The CSV file **must include a header row** with the columns in this order: `Content`, `Content Type`, `Policy`, `Tier`. Each row represents one blacklist entry.

  - **Content**: The text or regex pattern (see [Go regex syntax](https://pkg.go.dev/regexp/syntax)) to match against.
  - **Content Type**: `regex` for regular expressions, `string` for direct string matches
  - **Policy**: The name of the policy, which must correspond to a policy markdown file in the [config/policies](config/policies) directory.
  - **Tier**: `0` for words that trigger auto-rejection, `1` for words that indicate a possible violation but require further LLM evaluation.

### Policy Documents

- **Location and Naming**: Each policy referred to in the `Policy` column of the blacklist must have a corresponding markdown document in the [config/policies](config/policies) directory. For example, if a blacklist entry has the policy `profanity` then there must be a document describing that policy located here: [config/policies/profanity.md](config/policies/profanity.md)

- **Customization**: You are encouraged to review and modify the provided policy documents to fit your use-case. You can also create new policies by adding corresponding entries to the blacklist CSV and creating new policy markdown files.

## Interacting with the SMS Content Moderator API

The SMS Content Moderator is designed to function as an API that integrates into a messaging application. Below are the detailed specifications and examples of how the API processes and evaluates messages.

### Response Payload

| Field         | Type                                      | Description                                                                                                                    |
| ------------- | ----------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `status`      | string                                    | Overall status of the content evaluation. Possible values: `"pass"`, `"fail"`. Indicates whether any violations were detected. |
| `evaluations` | Array of [Evaluation](#evaluation-object) | A list of evaluation results for specific policies and tiers.                                                                  |

#### Evaluation Schema

Details about the evaluation of a specific piece of content against a defined policy and tier.

| Field       | Type   | Description                                                                                                                                                    |
| ----------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `status`    | string | The outcome of the evaluation for this policy. Values: `"is-violation"`, `"not-violation"`, `"not-evaluated"`.                                                 |
| `key`       | string | A composite key representing the tier and policy, formatted as `"{tier}-{policy}"`.                                                                            |
| `policy`    | string | The name of the policy that was evaluated.                                                                                                                     |
| `tier`      | int    | The severity tier of the violation (`0` or `1`). Tier 0 indicates severe violations leading to immediate rejection, while Tier 1 violations depend on context. |
| `reasoning` | string | Explanation of the evaluation outcome, providing context or justification for the result.                                                                      |

### Examples

#### Pass Without Violations

```json
{
  "status": "pass",
  "evaluations": []
}
```

### Fail with Tier 0 Violation

```json
{
  "status": "fail",
  "evaluations": [
    {
      "status": "is-violation",
      "key": "0-profanity",
      "policy": "profanity",
      "tier": 0,
      "reasoning": "Tier 0 blacklist entry was matched, which is automatically a policy violation."
    }
  ]
}
```

#### Pass with Tier 1 Blacklist Entry Matched

Tier 1 blacklist entries represent words that indicate a potential violation without

```json
{
  "status": "pass",
  "evaluations": [
    {
      "status": "not-violation",
      "key": "1-gambling",
      "policy": "gambling",
      "tier": 1,
      "reasoning": "The content is an invitation to a concert at a hotel and casino. It does not promote or facilitate gambling but mentions a casino as the venue for an event, which is a different context."
    }
  ]
}
```

### Fail with Tier 1 Violation

```json
{
  "status": "fail",
  "evaluations": [
    {
      "status": "is-violation",
      "key": "1-gambling",
      "policy": "gambling",
      "tier": 1,
      "reasoning": "The message encourages and promotes gambling activity at a casino, which is against the policy."
    }
  ]
}
```

### Fail with Tier 0 & Tier 1 Violations

Tier 0 blacklist entries signify a content policy has been violated. Tier 1 blacklist entries, which signify a policy _may_ have been violated, will not be evaluated after a Tier 0 violation.

```json
{
  "status": "fail",
  "evaluations": [
    {
      "status": "is-violation",
      "key": "0-profanity",
      "policy": "profanity",
      "tier": 0,
      "reasoning": "Tier 0 blacklist entry was matched, which is automatically a policy violation."
    },
    {
      "status": "not-evaluated",
      "key": "1-gambling",
      "policy": "gambling",
      "tier": 1,
      "reasoning": "Message content included a Tier 0 blacklist violation and there is no reason to evaluate Tier 1 policies."
    }
  ]
}
```
