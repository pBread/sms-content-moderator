package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/pbread/hoot-filter/internal/blacklist"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	initialize()

	http.HandleFunc("/webhook", handler)
	// http.HandleFunc("/webhook", auth.TwilioAuthMiddleware(handler))

	fmt.Println("Server is starting on port " + port + "...")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error starting server:", err)
		os.Exit(1)
	}
}

var (
	bl   *blacklist.Blacklist
	once sync.Once
)

func initialize() {
	once.Do(func() {
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
	})

	once.Do(func() {
		twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
		if len(twilioAuthToken) != 32 {
			panic("Invalid env variable: TWILIO_AUTH_TOKEN")
		}
	})
}

func handler(w http.ResponseWriter, req *http.Request) {
	initialize()

	fmt.Println("request")
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	msg := req.FormValue("Body")
	fmt.Println("body: " + msg)

	isBlacklistMatched := false

	isBlacklistMatched = bl.EvalTier0(msg)
	if isBlacklistMatched {
		http.Error(w, "Message contains a tier 0 prohibited word", http.StatusForbidden)
	}

	isBlacklistMatched = bl.EvalTier1(msg)
	if isBlacklistMatched {
		http.Error(w, "Message contains a tier 1 prohibited word", http.StatusForbidden)
	}

}
