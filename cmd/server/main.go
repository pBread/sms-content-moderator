package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/pBread/sms-content-moderator/internal/blacklist"
	"github.com/pBread/sms-content-moderator/internal/llm"
)

type RequestBody struct {
	Message string `json:"Message"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	projectRoot := getProjectRoot()
	csvPath := filepath.Join(projectRoot, "/config/blacklist.csv")
	blacklist.Init(csvPath)

	http.HandleFunc("/dev", unauthenticatedHandler)

	log.Println("starting on port" + ":8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func getProjectRoot() string {
	// retrieve the runtime file path
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Cannot retrieve runtime information")
	}

	// navigate up to the project root from current file (`cmd/server/main.go`)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(b)))

	return projectRoot
}

func unauthenticatedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Println()
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	violations := blacklist.CheckContent(reqBody.Message)

	prompt, err := llm.BuildPrompt(reqBody.Message, violations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(prompt)

	resp, _ := llm.EvalPolicyViolation(prompt)
	log.Println(resp)

	w.Header().Set("Content-Type", "application/json")
	response := map[string][]string{"Violations": violations}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}
