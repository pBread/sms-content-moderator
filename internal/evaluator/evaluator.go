package evaluator

import (
	"strconv"
	"strings"

	"github.com/pBread/sms-content-moderator/internal/blacklist"
	"github.com/pBread/sms-content-moderator/internal/llm"
)

type Evaluator interface {
	EvaluateContent(content string) (Response, error)
}

type ContentEvaluator struct{}

type Evaluation struct {
	Status    string `json:"status"`
	Key       string `json:"key"`
	Policy    string `json:"policy"`
	Tier      int    `json:"tier"`
	Reasoning string `json:"reasoning"`
}

type Response struct {
	Status      string       `json:"status"`
	Evaluations []Evaluation `json:"evaluations"`
}

// EvaluateContent checks the provided content against a set of blacklist rules and evaluates for policy violations using a tier-based system.
func (ce ContentEvaluator) EvaluateContent(content string) (Response, error) {
	// checks message for blacklist entry matches
	// returns []"{tier}-{policy}", e.g. ["0-profanity", "1-gambling"]
	blacklistMatches := blacklist.Match(content)

	result := Response{
		Status:      "pass",
		Evaluations: []Evaluation{},
	}

	tier0Present := false

	for _, match := range blacklistMatches {
		split := strings.Split(match, "-")
		tier, _ := strconv.Atoi(split[0])
		policy := split[1]

		if tier == 0 {
			result.Evaluations = append(result.Evaluations, Evaluation{
				Status:    "is-violation",
				Key:       match,
				Policy:    policy,
				Tier:      tier,
				Reasoning: "Tier 0 blacklist entry was matched, which is automatically a policy violation.",
			})
			result.Status = "fail"
			tier0Present = true
		}
	}

	if tier0Present {
		// generate 'not-evaluated' records for Tier 1 if Tier 0 is present
		for _, match := range blacklistMatches {
			split := strings.Split(match, "-")
			tier, _ := strconv.Atoi(split[0])
			if tier == 1 {
				result.Evaluations = append(result.Evaluations, Evaluation{
					Status:    "not-evaluated",
					Key:       match,
					Policy:    split[1],
					Tier:      tier,
					Reasoning: "Message content included a Tier 0 blacklist violation and there is no reason to evaluate Tier 1 policies.",
				})
			}
		}
	} else if result.Status == "pass" {
		// evaluate Tier 1 if no Tier 0 is present
		llmViolations, err := llm.AskLLM(content, blacklistMatches)
		if err != nil {
			return result, err
		}

		for _, violation := range llmViolations {
			result.Evaluations = append(result.Evaluations, Evaluation{
				Status:    violation.Status,
				Key:       "1-" + violation.Policy,
				Policy:    violation.Policy,
				Tier:      1,
				Reasoning: violation.Reasoning,
			})
			if violation.Status == "is-violation" {
				result.Status = "fail"
			}
		}
	}

	return result, nil
}
