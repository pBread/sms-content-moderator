package gpt

import (
	"log"

	"github.com/ayush6624/go-chatgpt"
)

var (
	client *chatgpt.Client
)

// returns true
func EvalContentPolicy(openAiKey string, msg string, topic string) bool {

}

func IsMsgOK(openAiKey string, msg string) {
	client, err := chatgpt.NewClient(openAiKey)
	if err != nil {
		log.Fatal(err)
	}

}
