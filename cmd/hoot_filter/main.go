package main

import (
	"net/http"
	"os"

	"github.com/pbread/hoot-filter/internal/auth"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	http.HandleFunc("/webhook", auth.TwilioAuthMiddleware(handler))
	http.ListenAndServe(port, nil)
}

func composeHandler() {

}

func handler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

}
