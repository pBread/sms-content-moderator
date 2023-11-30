package main

import (
	"fmt"
	"net/http"
	"sync"

	auth "github.com/pbread/hoot-filter/internal/auth"
	"github.com/pbread/hoot-filter/internal/blacklist"
)

// func main() {
// 	http.HandleFunc("/webhook", handler)
// 	http.ListenAndServe(":8080", nil)
// }

func main() {
	bl := blacklist.GetBlackList()

	fmt.Printf("bl: %v\n", bl)

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

	if err := req.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var tier0Failed bool
	var tier1Failed bool

	msg := req.FormValue("Body")

	go func() {
		defer wg.Done()
		checkTier0(&mutex, msg, &tier0Failed)
	}()

	go func() {
		defer wg.Done()
		checkTier1(&mutex, msg, &tier0Failed, &tier1Failed)
	}()
}

func checkTier0(mutex *sync.Mutex, msg string, tier0Failed *bool) {
}

func checkTier1(mutex *sync.Mutex, msg string, tier0Failed *bool, tier1Failed *bool) {
}

func containsWord(msg string, word string) {}
