/****************************************************
 Only Supports POST webhooks
****************************************************/

package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const twilioAuthToken = "TWILIO_AUTH_TOKEN"

func TwilioAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isSignatureValid, err := validateRequest(r)
		if err != nil {
			fmt.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if !isSignatureValid {
			w.WriteHeader(http.StatusUnauthorized)
		}

		next(w, r)
	}
}

func validateRequest(req *http.Request) (bool, error) {
	signatureGiven := req.Header.Get("X-Twilio-Signature")

	signatureExpected, err := sign(req)
	if err != nil {
		return false, err
	}

	// Compare the calculated signature with the expected signature
	return hmac.Equal([]byte(signatureGiven), []byte(signatureExpected)), nil
}

func sign(req *http.Request) (string, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", errors.New("No request body. Hint: This app only supports POST requests.")
	}

	// Create the base string to sign, which is the full URL of the request
	// concatenated with the URL-encoded POST body
	url := "https://" + req.Host + req.URL.String()
	dataToSign := url + string(body)

	// Create an HMAC-SHA1 hasher
	hasher := hmac.New(sha1.New, []byte(twilioAuthToken))
	// Write the data to the hasher
	hasher.Write([]byte(dataToSign))

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
