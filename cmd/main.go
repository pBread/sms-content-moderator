package main

import (
	"fmt"
	"net/http"

	auth "github.com/pbread/hoot-filter/internal/auth"
)

func main() {
	http.HandleFunc("/webhook", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, req *http.Request) {
	isSignatureValid, err := auth.ValidateTwilioRequest(req)
	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if !isSignatureValid {
		w.WriteHeader(http.StatusUnauthorized)
	}

}
