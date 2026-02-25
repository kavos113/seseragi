package golang

import (
	"encoding/json"
	"os"
)

type Empty struct{}

type NodeData struct {
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

type InputData []NodeData

func (d InputData) Get(name string, out *any) error {
	for _, node := range d {
		if node.Name == name {
			return json.Unmarshal(node.Data, out)
		}
	}
	return nil
}

func Run[O any](handler func(data InputData) (O, error)) {
	inPath := os.Getenv("WORKFLOW_INPUT_PATH")
	outPath := os.Getenv("WORKFLOW_OUTPUT_PATH")

	inData, err := os.ReadFile(inPath)
	if err != nil {
		writeError(outPath, err)
		return
	}

	var input InputData
	if err := json.Unmarshal(inData, &input); err != nil {
		writeError(outPath, err)
		return
	}

	outData, err := handler(input)
	if err != nil {
		writeError(outPath, err)
		return
	}

	outBytes, err := json.Marshal(outData)
	if err != nil {
		writeError(outPath, err)
		return
	}

	err = os.WriteFile(outPath, outBytes, 0644)
	if err != nil {
		writeError(outPath, err)
	}
}

func writeError(path string, err error) {
	errData, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	os.WriteFile(path, errData, 0644)
}
