package main

import seseragi "github.com/kavos113/seseragi/sdk/golang"

type Message struct {
	Question string `json:"question"`
}

func askName(_ seseragi.Empty) (Message, error) {
	return Message{
		Question: "What is your name?",
	}, nil
}

func main() {
	seseragi.Run(askName)
}
