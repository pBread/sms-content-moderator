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

func TwilioAuthMiddleware(next http.HandlerFunc, authToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		signatureGiven := req.Header.Get("X-Twilio-Signature")

		signatureExpected, err := sign(req, authToken)
		if err != nil {
			fmt.Println("Error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		isSignatureValid := hmac.Equal([]byte(signatureGiven), []byte(signatureExpected))
		if err != nil {
			fmt.Println("Error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if !isSignatureValid {
			w.WriteHeader(http.StatusUnauthorized)
		}

		next(w, req)
	}
}

func sign(req *http.Request, authToken string) (string, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", errors.New("No request body. Hint: This app only supports POST requests.")
	}

	// Create the base string to sign, which is the full URL of the request
	// concatenated with the URL-encoded POST body
	url := "https://" + req.Host + req.URL.String()
	dataToSign := url + string(body)

	// Create an HMAC-SHA1 hasher
	hasher := hmac.New(sha1.New, []byte(authToken))
	// Write the data to the hasher
	hasher.Write([]byte(dataToSign))

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
