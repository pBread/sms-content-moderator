package main

import (
	"encoding/json"
	"net/http"

	"github.com/pBread/sms-content-moderator/internal/evaluator"
	"github.com/pBread/sms-content-moderator/internal/logger"
)

type RequestBody struct {
	Message string `json:"Message"`
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

	response, err := evaluator.EvaluateContent(reqBody.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Error writing response: ", err.Error())
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}
