package manager

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/kavos113/seseragi/model"
)

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

	for nodeName, nodeInfo := range workflowInfo.Nodes {
		currentNodeIndex := slices.IndexFunc(nodes, func(n model.Node) bool {
			return n.TaskID == nodeInfo.ID
		})
		if currentNodeIndex == -1 {
			return model.Workflow{}, fmt.Errorf("node %s not found in nodes list", nodeName)
		}
		currentNode := &nodes[currentNodeIndex]

		for _, depName := range nodeInfo.Dependencies {
			depNodeInfo, ok := workflowInfo.Nodes[depName]
			if !ok {
				return model.Workflow{}, fmt.Errorf("dependency %s not found for node %s", depName, nodeName)
			}

			depNodeIndex := slices.IndexFunc(nodes, func(n model.Node) bool {
				return n.TaskID == depNodeInfo.ID
			})
			if depNodeIndex == -1 {
				return model.Workflow{}, fmt.Errorf("dependency node %s not found in nodes list", depNodeInfo.ID)
			}

			currentNode.Dependencies = append(currentNode.Dependencies, &nodes[depNodeIndex])
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
