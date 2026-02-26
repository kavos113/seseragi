package yaml

import (
	"testing"

	"github.com/kavos113/seseragi/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestLoadWorkflowInfoFromYAML(t *testing.T) {
	tests := []struct {
		name     string
		yamlData []byte
		yamlPath string
		want     *domain.Workflow
		wantErr  bool
	}{
		{
			name: "valid YAML",
			yamlData: []byte(`
name: "hello-workflow"
description: "Hello Workflow"
nodes:
  go-hello:
    name: "go-hello"
`),
			yamlPath: "flow.yaml",
			want: &domain.Workflow{
				Name: "hello-workflow",
				Nodes: []domain.Node{
					{
						Name:     "go-hello",
						TaskName: "go-hello",
					},
				},
				YamlPath: "flow.yaml",
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
    name: "go-hello"

  go-world:
    name: "go-world"
    dependencies:
      - go-hello
    environments:
      VAR1: value1
`),
			yamlPath: "flow.yaml",
			want: &domain.Workflow{
				Name: "hello-workflow",
				Nodes: []domain.Node{
					{
						Name:     "go-hello",
						TaskName: "go-hello",
					},
					{
						Name:         "go-world",
						TaskName:     "go-world",
						Dependencies: []string{"go-hello"},
						Environments: map[string]string{"VAR1": "value1"},
					},
				},
				YamlPath: "flow.yaml",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadWorkflowInfoFromYAML(tt.yamlData, tt.yamlPath)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("LoadWorkflowInfoFromYAML() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.YamlPath, got.YamlPath)
			assert.Equal(t, len(tt.want.Nodes), len(got.Nodes))
			for i, node := range tt.want.Nodes {
				assert.Equal(t, node.Name, got.Nodes[i].Name)
				assert.Equal(t, node.TaskName, got.Nodes[i].TaskName)
				assert.Equal(t, node.Dependencies, got.Nodes[i].Dependencies)
			}
		})
	}
}
