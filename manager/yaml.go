package manager

import "go.yaml.in/yaml/v4"

type TaskInfo struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Context     string `yaml:"context"`
	Path        string
}

func LoadTaskInfoFromYAML(yamlData []byte, yamlPath string) (*TaskInfo, error) {
	var taskInfo TaskInfo
	err := yaml.Unmarshal(yamlData, &taskInfo)
	if err != nil {
		return nil, err
	}
	taskInfo.Path = yamlPath
	return &taskInfo, nil
}
