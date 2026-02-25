package yaml

import (
	"path/filepath"
	"testing"

	"github.com/kavos113/seseragi/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestLoadTaskInfoFromYAML(t *testing.T) {
	tests := []struct {
		name     string
		yamlData []byte
		yamlPath string
		want     *domain.Task
		wantErr  bool
	}{
		{
			name: "valid YAML",
			yamlData: []byte(`
name: "go-hello"
description: "Go Hello Task"
type: "docker"
context: .
`),
			yamlPath: "task.yaml",
			want: &domain.Task{
				Name:     "go-hello",
				YamlPath: "task.yaml",
				TaskDef: domain.DockerTaskDefinition{
					ContextDir: ".",
				},
			},
			wantErr: false,
		},
		{
			name: "yaml with path",
			yamlData: []byte(`
name: "go-hello"
description: "Go Hello Task"
type: "docker"
context: .
path: "task.yaml"
`),
			yamlPath: "other.yaml",
			want: &domain.Task{
				Name:     "go-hello",
				YamlPath: "other.yaml",
				TaskDef: domain.DockerTaskDefinition{
					ContextDir: ".",
				},
			},
			wantErr: false,
		},
		{
			name: "command type",
			yamlData: []byte(`
name: "go-hello"
description: "Go Hello Task"
type: "command"
command: "echo Hello"
working_dir: .
`),
			yamlPath: "task.yaml",
			want: &domain.Task{
				Name:     "go-hello",
				YamlPath: "task.yaml",
				TaskDef: domain.CommandTaskDefinition{
					Command:    "echo Hello",
					WorkingDir: ".",
				},
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
			assert.Equal(t, got.YamlPath, tt.want.YamlPath)

			switch got.TaskDef.Type() {
			case domain.TaskTypeDocker:
				gotDef := got.TaskDef.(domain.DockerTaskDefinition)
				wantDef := tt.want.TaskDef.(domain.DockerTaskDefinition)

				absYaml, err := filepath.Abs(tt.yamlPath)
				if err != nil {
					t.Errorf("Failed to get absolute path: %v", err)
				}
				assert.Equal(t, gotDef.ContextDir, filepath.Join(filepath.Dir(absYaml), wantDef.ContextDir))

			case domain.TaskTypeCommand:
				gotDef := got.TaskDef.(domain.CommandTaskDefinition)
				wantDef := tt.want.TaskDef.(domain.CommandTaskDefinition)
				assert.Equal(t, gotDef.Command, wantDef.Command)
				assert.Equal(t, gotDef.WorkingDir, wantDef.WorkingDir)
			}
		})
	}
}
