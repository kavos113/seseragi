package yaml

import (
	"os"
	"regexp"
	"time"

	"github.com/kavos113/seseragi/internal/domain"
	"go.yaml.in/yaml/v4"
)

var (
	envVarRegex = regexp.MustCompile(`\$\{\s*([A-Za-z_][A-Za-z0-9_]*)\s*\}`)
)

type nodeInfo struct {
	Name         string            `yaml:"name"`
	Dependencies []string          `yaml:"dependencies"`
	Environments map[string]string `yaml:"environments,omitempty"`
}

type triggerInfo struct {
	Type     string `yaml:"type"`
	Interval int    `yaml:"interval,omitempty"` // minutes
}

type workflowInfo struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Nodes       map[string]nodeInfo `yaml:"nodes"`
	Trigger     triggerInfo         `yaml:"trigger"`
}

// return ID, is empty
func LoadWorkflowInfoFromYAML(yamlData []byte, yamlPath string) (*domain.Workflow, error) {
	// environment variables
	expandedYamlData := envVarRegex.ReplaceAllFunc(yamlData, func(match []byte) []byte {
		envVar := envVarRegex.FindSubmatch(match)[1]
		val := os.Getenv(string(envVar))
		if val == "" {
			return match // keep original if env var is not set
		}
		return []byte(val)
	})

	var workflowInfo workflowInfo
	err := yaml.Unmarshal(expandedYamlData, &workflowInfo)
	if err != nil {
		return nil, err
	}

	nodes := make([]domain.Node, 0, len(workflowInfo.Nodes))
	for name, node := range workflowInfo.Nodes {
		nodes = append(nodes, domain.Node{
			Name:         name,
			TaskName:     node.Name,
			Dependencies: node.Dependencies,
			Environments: node.Environments,
		})
	}

	workflow := &domain.Workflow{
		Name:     workflowInfo.Name,
		Nodes:    nodes,
		YamlPath: yamlPath,
	}

	if workflowInfo.Trigger.Type == "interval" {
		workflow.RunInterval = time.Duration(workflowInfo.Trigger.Interval) * time.Minute
	} else {
		workflow.RunInterval = time.Duration(time.Hour * 24 * 365 * 100) // effectively never run
	}

	return workflow, nil
}
