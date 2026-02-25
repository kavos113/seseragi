package golang

import (
	"encoding/json"
	"os"
)

type Empty struct{}

func Run[I any, O any](handler func(I) (O, error)) {
	inPath := os.Getenv("WORKFLOW_INPUT_PATH")
	outPath := os.Getenv("WORKFLOW_OUTPUT_PATH")

	inData, err := os.ReadFile(inPath)
	if err != nil {
		writeError(outPath, err)
		os.Exit(1)
	}

	var input I
	err = json.Unmarshal(inData, &input)
	if err != nil {
		writeError(outPath, err)
		os.Exit(1)
	}

	output, err := handler(input)
	if err != nil {
		writeError(outPath, err)
		os.Exit(1)
	}

	outData, err := json.Marshal(output)
	if err != nil {
		writeError(outPath, err)
		os.Exit(1)
	}

	err = os.WriteFile(outPath, outData, 0644)
	if err != nil {
		writeError(outPath, err)
		os.Exit(1)
	}
}

func writeError(path string, err error) {
	errData, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	os.WriteFile(path, errData, 0644)
}