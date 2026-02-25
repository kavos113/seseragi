package yaml

import (
	"github.com/kavos113/seseragi/internal/domain"
	"go.yaml.in/yaml/v4"
)

type nodeInfo struct {
	Name         string   `yaml:"name"`
	Dependencies []string `yaml:"dependencies"`
}

type workflowInfo struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Nodes       map[string]nodeInfo `yaml:"nodes"`
}

// return ID, node.TaskID is empty
func LoadWorkflowInfoFromYAML(yamlData []byte, yamlPath string) (*domain.Workflow, error) {
	var workflowInfo workflowInfo
	err := yaml.Unmarshal(yamlData, &workflowInfo)
	if err != nil {
		return nil, err
	}

	nodes := make([]domain.Node, 0, len(workflowInfo.Nodes))
	for name, node := range workflowInfo.Nodes {
		nodes = append(nodes, domain.Node{
			Name:         name,
			TaskName:     node.Name,
			Dependencies: node.Dependencies,
		})
	}

	return &domain.Workflow{
		Name:     workflowInfo.Name,
		Nodes:    nodes,
		YamlPath: yamlPath,
	}, nil
}
