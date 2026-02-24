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
		{
			name: "valid YAML with dependency",
			yamlData: []byte(`
name: "hello-workflow"
description: "Hello Workflow"

nodes:
  go-hello:
    id: "some-id"

  go-world:
    id: "some-other-id"
    dependencies:
      - go-hello
`),
			yamlPath: "flow.yaml",
			want: &WorkflowInfo{
				Name:        "hello-workflow",
				Description: "Hello Workflow",
				Nodes: map[string]NodeInfo{
					"go-hello": {ID: "some-id"},
					"go-world": {ID: "some-other-id", Dependencies: []string{"go-hello"}},
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

			for key, node := range tt.want.Nodes {
				gotNode, exists := got.Nodes[key]
				assert.True(t, exists, "Node %s should exist", key)
				assert.Equal(t, gotNode.ID, node.ID)
				assert.Equal(t, gotNode.Dependencies, node.Dependencies)
			}
		})
	}
}
