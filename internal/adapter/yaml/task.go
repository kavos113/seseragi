package yaml

import (
	"errors"
	"path/filepath"

	"github.com/kavos113/seseragi/internal/domain"
	"go.yaml.in/yaml/v4"
)

type taskInfo struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
}

type dockerTaskDetails struct {
	Context string `yaml:"context"`
}

type commandTaskDetails struct {
	Command    string `yaml:"command"`
	WorkingDir string `yaml:"working_dir"`
}

// return ID is empty
func LoadTaskInfoFromYAML(yamlData []byte, yamlPath string) (*domain.Task, error) {
	var taskInfo taskInfo
	err := yaml.Unmarshal(yamlData, &taskInfo)
	if err != nil {
		return nil, err
	}

	switch taskInfo.Type {
	case "docker":
		var details dockerTaskDetails
		err = yaml.Unmarshal(yamlData, &details)
		if err != nil {
			return nil, err
		}

		absYaml, err := filepath.Abs(yamlPath)
		if err != nil {
			return nil, err
		}
		contextDir := filepath.Join(filepath.Dir(absYaml), details.Context)

		return &domain.Task{
			Name:     taskInfo.Name,
			YamlPath: yamlPath,
			TaskDef: domain.DockerTaskDefinition{
				ContextDir: contextDir,
			},
		}, nil

	case "command":
		var details commandTaskDetails
		err = yaml.Unmarshal(yamlData, &details)
		if err != nil {
			return nil, err
		}
		return &domain.Task{
			Name:     taskInfo.Name,
			YamlPath: yamlPath,
			TaskDef: domain.CommandTaskDefinition{
				Command:    details.Command,
				WorkingDir: details.WorkingDir,
			},
		}, nil

	default:
		return nil, errors.New("unsupported task type: " + taskInfo.Type)
	}
}
