package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadTastInfoFromYAML(t *testing.T) {
	tests := []struct {
		name     string
		yamlData []byte
		yamlPath string
		want     *TaskInfo
		wantErr  bool
	}{
		{
			name: "valid YAML",
			yamlData: []byte(`
name: "go-hello"
description: "Go Hello Task"
context: .
`),
			yamlPath: "task.yaml",
			want: &TaskInfo{
				Name:        "go-hello",
				Description: "Go Hello Task",
				Context:     ".",
				Path:        "task.yaml",
			},
			wantErr: false,
		},
		{
			name: "yaml with path",
			yamlData: []byte(`
name: "go-hello"
description: "Go Hello Task"
context: .
path: "task.yaml"
`),
			yamlPath: "other.yaml",
			want: &TaskInfo{
				Name:        "go-hello",
				Description: "Go Hello Task",
				Context:     ".",
				Path:        "other.yaml",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadTaskInfoFromYAML(tt.yamlData, tt.yamlPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTaskInfoFromYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got.Name, tt.want.Name)
			assert.Equal(t, got.Description, tt.want.Description)
			assert.Equal(t, got.Context, tt.want.Context)
			assert.Equal(t, got.Path, tt.want.Path)
		})
	}
}

func TestLoadWorkflowInfoFromYAML(t *testing.T) {
	tests := []struct {
		name     string
		yamlData []byte
		yamlPath string
		want     *WorkflowInfo
		wantErr  bool
	}{
		{
			name: "valid YAML",
			yamlData: []byte(`
name: "hello-workflow"
description: "Hello Workflow"
nodes:
  go-hello:
    id: "some-id"
`),
			yamlPath: "flow.yaml",
			want: &WorkflowInfo{
				Name:        "hello-workflow",
				Description: "Hello Workflow",
				Nodes: map[string]NodeInfo{
					"go-hello": {ID: "some-id"},
				},
				Path: "flow.yaml",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadWorkflowInfoFromYAML(tt.yamlData, tt.yamlPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadWorkflowInfoFromYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got.Name, tt.want.Name)
			assert.Equal(t, got.Description, tt.want.Description)
			assert.Equal(t, got.Path, tt.want.Path)
			assert.Equal(t, len(got.Nodes), len(tt.want.Nodes))
			assert.Equal(t, got.Nodes["go-hello"].ID, tt.want.Nodes["go-hello"].ID)
		})
	}
}
