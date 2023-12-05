package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pbread/hoot-filter/internal/blacklist"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	prepare()

	http.HandleFunc("/webhook", handler)
	// http.HandleFunc("/webhook", auth.TwilioAuthMiddleware(handler))

	fmt.Println("Server is starting on port " + port + "...")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error starting server:", err)
		os.Exit(1)
	}

}

func prepare() {
	// prepares blacklist
	blacklist.GetBlackList()
}

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("request")

	if err := req.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	msg := req.FormValue("Body")
	fmt.Println("body: " + msg)

	bl := blacklist.GetBlackList()
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
