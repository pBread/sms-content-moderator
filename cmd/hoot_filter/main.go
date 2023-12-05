package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pbread/hoot-filter/internal/blacklist"
)

var tier0Fails = [...]string{
	"very bad message",
	"this message is very bad",
	"  very bad   with whitespace",
	"v e r y b a d",
	"badregex",
	"short message that contains bad regex",
	"short message that contains bad regex at the end",
	"long message with violation at end. Some Characters BAM BADRenter REGEXBAD. Let$s go. ab abc abcd abcde abcdef abcdefg acbdefghijklm acbdefghijklmopqrstuvwxyz acbdef ghijklmo pqrstuvwxyz. ab abc abcd abcde abcdef abcdefg acbdefghijklm acbdefghijklmopqrstuvwxyz acbdef ghijklmo pqrstuvwxyz. ab abc abcd abcde abcdef abcdefg acbdefghijklm acbdefghijklmopqrstuvwxyz acbdef ghijklmo pqrstuvwxyz. ab abc abcd abcde abcdef abcdefg acbdefghijklm acbdefghijklmopqrstuvwxyz acbdef ghijklmo pqrstuvwxyz. ab abc abcd abcde abcdef abcdefg acbdefghijklm acbdefghijklmopqrstuvwxyz acbdef ghijklmo pqrstuvwxyz. ab abc abcd abcde abcdef abcdefg acbdefghijklm acbdefghijklmopqrstuvwxyz acbdef ghijklmo pqrstuvwxyz. ab abc abcd abcde abcdef abcdefg acbdefghijklm acbdefghijklmopqrstuvwxyz acbdef ghijklmo pqrstuvwxyz.  that contains bad regex.",
	"should pass",
	"should pass",
	"should pass",
	"should pass",
	"should pass",
	"should pass",
	"should pass",
	"should pass",
	"should pass",
	"should pass",
	"should pass",
	"should pass",
	"should pass",
}

func main() {
	bl := blacklist.GetBlackList()

	for idx, msg := range tier0Fails {
		syncEvalTier0(bl, idx, msg)
	}
}

func syncEvalTier0(bl *blacklist.Blacklist, idx int, msg string) {
	const runs = 1000
	var totalDuration time.Duration

	for i := 0; i < runs; i++ {
		startTime := time.Now()
		_ = bl.SyncCheckTier0(msg)
		totalDuration += time.Since(startTime)
	}

	avgDuration := totalDuration.Nanoseconds() / runs
	result := bl.SyncCheckTier0(msg)

	resultStr := "failed"
	if result {
		resultStr = "passed\t"
	}
	fmt.Printf("Msg\t %d\t %q\t Avg time: %v\t Result: %s\n", idx, firstN(msg, 25), avgDuration, resultStr)
}

func firstN(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	}
	return s + strings.Repeat(" ", n-len(s))
}

// func main() {
// 	http.HandleFunc("/webhook", handler)
// 	http.ListenAndServe(":8080", nil)
// }

// func handler(w http.ResponseWriter, req *http.Request) {
// 	isSignatureValid, err := auth.ValidateTwilioRequest(req)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 	}
// 	if !isSignatureValid {
// 		w.WriteHeader(http.StatusUnauthorized)
// 	}

// 	if err := req.ParseForm(); err != nil {
// 		http.Error(w, "Error parsing form", http.StatusBadRequest)
// 		return
// 	}

// }
