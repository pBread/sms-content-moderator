package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pBread/sms-content-moderator/internal/blacklist"
	"github.com/pBread/sms-content-moderator/internal/logger"
)

type RequestBody struct {
	Message string `json:"Message"`
}

type Evaluation struct {
	Status    string `json:"status"`
	Key       string `json:"key"`
	Policy    string `json:"policy"`
	Tier      int    `json:"tier"`
	Reasoning string `json:"reasoning"`
}

type Response struct {
	Status      string       `json:"status"`
	Matches     []string     `json:"matches"`
	Evaluations []Evaluation `json:"evaluations"`
}

func main() {
	http.HandleFunc("/evaluate-message", unauthenticatedHandler)
	logger.Info("Starting on port" + ":8080")
	logger.Fatal(http.ListenAndServe(":8080", nil))
}

func unauthenticatedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logger.Error("Error reading request body: ", err.Error())
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	blacklistMatches := blacklist.CheckContent(reqBody.Message)
	evaluations := []Evaluation{}
	overallStatus := "pass"

	tier0Present := false

	for _, match := range blacklistMatches {
		split := strings.Split(match, "-")
		tier, _ := strconv.Atoi(split[0])
		policy := split[1]

		if tier == 0 {
			evaluations = append(evaluations, Evaluation{
				Status:    "is-violation",
				Key:       match,
				Policy:    policy,
				Tier:      tier,
				Reasoning: "Tier 0 blacklist entry matched, which is automatically a policy violation.",
			})
			overallStatus = "fail"
			tier0Present = true
		}
	}

	if tier0Present { // generate 'not-evaluated' records for Tier 1 if Tier 0 is present

	} else if overallStatus == "pass" { // evaluate Tier 1 if no Tier 0 is present

	}

	response := Response{
		Status:      overallStatus,
		Matches:     blacklistMatches,
		Evaluations: evaluations,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Error writing response: ", err.Error())
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}
