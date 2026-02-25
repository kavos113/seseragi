package main

import (
	"fmt"

	seseragi "github.com/kavos113/seseragi/sdk/golang"
)

type Message struct {
	Question string `json:"question"`
}

func askName(message Message) (seseragi.Empty, error) {
	fmt.Printf("Question: %s\n", message.Question)

	return seseragi.Empty{}, nil
}

func main() {
	seseragi.Run(askName)
}
