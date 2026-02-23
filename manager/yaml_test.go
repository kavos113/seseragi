package manager

import "testing"

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

			if got.Name != tt.want.Name {
				t.Errorf("LoadTaskInfoFromYAML() got.Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Description != tt.want.Description {
				t.Errorf("LoadTaskInfoFromYAML() got.Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.Context != tt.want.Context {
				t.Errorf("LoadTaskInfoFromYAML() got.Context = %v, want %v", got.Context, tt.want.Context)
			}
			if got.Path != tt.want.Path {
				t.Errorf("LoadTaskInfoFromYAML() got.Path = %v, want %v", got.Path, tt.want.Path)
			}
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

			if got.Name != tt.want.Name {
				t.Errorf("LoadWorkflowInfoFromYAML() got.Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Description != tt.want.Description {
				t.Errorf("LoadWorkflowInfoFromYAML() got.Description = %v, want %v", got.Description, tt.want.Description)
			}
			if len(got.Nodes) != len(tt.want.Nodes) {
				t.Errorf("LoadWorkflowInfoFromYAML() got.Nodes length = %v, want %v", len(got.Nodes), len(tt.want.Nodes))
			} else {
				for key, gotNode := range got.Nodes {
					wantNode, exists := tt.want.Nodes[key]
					if !exists {
						t.Errorf("LoadWorkflowInfoFromYAML() got.Nodes has unexpected key = %v", key)
						continue
					}
					if gotNode.ID != wantNode.ID {
						t.Errorf("LoadWorkflowInfoFromYAML() got.Nodes[%v].ID = %v, want %v", key, gotNode.ID, wantNode.ID)
					}
				}
			}
			if got.Path != tt.want.Path {
				t.Errorf("LoadWorkflowInfoFromYAML() got.Path = %v, want %v", got.Path, tt.want.Path)
			}
		})
	}
}