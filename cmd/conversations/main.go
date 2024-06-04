package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go/client"
)

var twilioAuthToken string

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Unable to load environment variables")
	}
	twilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	if twilioAuthToken == "" {
		log.Fatal("Missing env variable: TWILIO_AUTH_TOKEN")
	}

	http.HandleFunc("/conversations-pre-event", conversationsPreEventHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func conversationsPreEventHandler(w http.ResponseWriter, r *http.Request) {
	// validate twilio signature
	conversationsWebhookUrl := os.Getenv("CONVERSATIONS_PRE_EVENT_WEBHOOK_URL")
	if conversationsWebhookUrl == "" {
		log.Fatal("Missing env variable: CONVERSATIONS_PRE_EVENT_WEBHOOK_URL")
	}

	bodyBytes, err := io.ReadAll(r.Body) // Read the body into a byte slice
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	requestValidator := client.NewRequestValidator(twilioAuthToken)
	requestValidator.ValidateBody(conversationsWebhookUrl, bodyBytes, r.Header["X-Twilio-Signature"][0])

}
