package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pbread/hoot-filter/internal/auth"
	"github.com/pbread/hoot-filter/internal/blacklist"
)

var (
	openAiKey       string
	twilioAuthToken string

	bl *blacklist.Blacklist
)

func main() {
	loadEnv()
	loadBlacklist()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/webhook", handler)
	http.HandleFunc("/webhook-with-auth", auth.TwilioAuthMiddleware(handler, twilioAuthToken))

	fmt.Println("Server is starting on port " + port + "...")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error starting server:", err)
		os.Exit(1)
	}
}

func loadEnv() {
	godotenv.Load()

	openAiKey = os.Getenv("OPENAI_KEY")
	if len(twilioAuthToken) == 0 {
		panic("Missing env variable: OPENAI_KEY")
	}

	twilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	if len(twilioAuthToken) != 32 {
		panic("Missing or invalid env variable: TWILIO_AUTH_TOKEN")
	}
}

func loadBlacklist() {
	file, err := os.Open("config/blacklist.csv")
	if err != nil {
		panic("Error loading blacklist CSV: " + err.Error())
	}

	reader := csv.NewReader(file)
	entries, err := reader.ReadAll()
	if err != nil {
		panic("Error reading blacklist CSV: " + err.Error())
	}

	bl = blacklist.MakeBlacklist(entries)
	fmt.Println("blacklist initialized")

}

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("request")
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	msg := req.FormValue("Body")
	fmt.Println("body: " + msg)

	isBlacklistMatched := bl.EvalTier0(msg)
	if isBlacklistMatched {
		http.Error(w, "Message contains a tier 0 prohibited word", http.StatusForbidden)
	}

	isBlacklistMatched = bl.EvalTier1(msg)
	if isBlacklistMatched {
		http.Error(w, "Message contains a tier 1 prohibited word", http.StatusForbidden)
	}
}
