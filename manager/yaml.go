package manager

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/kavos113/seseragi/model"
	"go.yaml.in/yaml/v4"
)

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

type NodeInfo struct {
	ID           string   `yaml:"id"`
	Dependencies []string `yaml:"dependencies"`
}

type WorkflowInfo struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Nodes       map[string]NodeInfo `yaml:"nodes"`
	Path        string
}

func LoadWorkflowInfoFromYAML(yamlData []byte, yamlPath string) (*WorkflowInfo, error) {
	var workflowInfo WorkflowInfo
	err := yaml.Unmarshal(yamlData, &workflowInfo)
	if err != nil {
		return nil, err
	}
	workflowInfo.Path = yamlPath
	return &workflowInfo, nil
}

func ParseWorkflow(workflowInfo *WorkflowInfo, yamlPath string) (model.Workflow, error) {
	nodes := make([]model.Node, 0, len(workflowInfo.Nodes))
	for nodeName, nodeInfo := range workflowInfo.Nodes {
		task, err := taskRepo.GetTaskByID(nodeInfo.ID)
		if err != nil {
			return model.Workflow{}, fmt.Errorf("failed to get task by ID %s for node %s: %w", nodeInfo.ID, nodeName, err)
		}
		nodes = append(nodes, model.Node{
			TaskID:       task.ID,
			Dependencies: []*model.Node{},
		})
	}

	for _, node := range nodes {
		dependencies := workflowInfo.Nodes[node.TaskID].Dependencies
		for _, depName := range dependencies {
			depNodeInfo, ok := workflowInfo.Nodes[depName]
			if !ok {
				return model.Workflow{}, fmt.Errorf("dependency %s not found for node %s", depName, node.TaskID)
			}
			depNodeIndex := slices.IndexFunc(nodes, func(n model.Node) bool {
				return n.TaskID == depNodeInfo.ID
			})
			if depNodeIndex == -1 {
				return model.Workflow{}, fmt.Errorf("dependency node %s not found in nodes list", depNodeInfo.ID)
			}
			node.Dependencies = append(node.Dependencies, &nodes[depNodeIndex])
		}
	}

	id := uuid.New().String()
	workflow := model.Workflow{
		ID:       id,
		Name:     workflowInfo.Name,
		Nodes:    nodes,
		YamlPath: yamlPath,
	}
	return workflow, nil
}
