package golang

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputData_Get(t *testing.T) {
	data := InputData{
		{Name: "node1", Data: json.RawMessage(`"value1"`)},
		{Name: "node2", Data: json.RawMessage(`{"key": "value2"}`)},
	}

	var s string
	err := data.Get("node1", &s)
	assert.NoError(t, err)
	assert.Equal(t, "value1", s)

	var m map[string]string
	err = data.Get("node2", &m)
	assert.NoError(t, err)
	assert.Equal(t, "value2", m["key"])

	var notFound string
	err = data.Get("node3", &notFound)
	assert.NoError(t, err)
	assert.Equal(t, "", notFound)
}

func TestRun(t *testing.T) {
	inPath := t.TempDir() + "/input.json"
	outPath := t.TempDir() + "/output.json"

	input := InputData{
		{Name: "node1", Data: json.RawMessage(`"hello"`)},
	}
	inputBytes, _ := json.Marshal(input)
	os.WriteFile(inPath, inputBytes, 0644)

	os.Setenv("WORKFLOW_INPUT_PATH", inPath)
	os.Setenv("WORKFLOW_OUTPUT_PATH", outPath)
	defer os.Unsetenv("WORKFLOW_INPUT_PATH")
	defer os.Unsetenv("WORKFLOW_OUTPUT_PATH")

	handler := func(data InputData) (string, error) {
		var s string
		data.Get("node1", &s)
		return s + " world", nil
	}

	Run(handler)

	outBytes, err := os.ReadFile(outPath)
	assert.NoError(t, err)

	var output string
	err = json.Unmarshal(outBytes, &output)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", output)
}
