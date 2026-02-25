package main

import (
	"fmt"

	seseragi "github.com/kavos113/seseragi/sdk/golang"
)

type Message struct {
	Question string `json:"question"`
}

func askName(d seseragi.InputData) (seseragi.Empty, error) {
	var message Message
	if err := d.Get("question", &message); err != nil {
		return seseragi.Empty{}, err
	}
	fmt.Printf("Question: %s\n", message.Question)

	return seseragi.Empty{}, nil
}

func main() {
	seseragi.Run(askName)
}
