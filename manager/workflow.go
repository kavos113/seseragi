package manager

import (
	"fmt"
	"slices"
	"strings"

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
			Dependencies: []string{},
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

			currentNode.Dependencies = append(currentNode.Dependencies, depNodeInfo.ID)
		}
	}

	if err := checkCircularDependency(nodes); err != nil {
		return model.Workflow{}, err
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

func checkCircularDependency(nodes []model.Node) error {
	// 0: unvisit, 1: visiting, 2: visited
	state := make(map[string]int)

	var visit func(nodeId string, stack []string) error
	visit = func(nodeId string, stack []string) error {
		if state[nodeId] == 1 {
			return fmt.Errorf("%w: %s", ErrWorkflowCircularDependency, strings.Join(append(stack, nodeId), " -> "))
		}
		if state[nodeId] == 2 {
			return nil
		}
		state[nodeId] = 1

		for _, dep := range getDependencies(nodes, nodeId) {
			if err := visit(dep, append(stack, dep)); err != nil {
				return err
			}
		}
		state[nodeId] = 2
		return nil
	}

	for _, node := range nodes {
		if err := visit(node.TaskID, []string{}); err != nil {
			return err
		}
	}
	return nil
}

func getDependencies(nodes []model.Node, nodeId string) []string {
	for _, node := range nodes {
		if node.TaskID == nodeId {
			return node.Dependencies
		}
	}
	return []string{}
}
