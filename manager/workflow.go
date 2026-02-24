package manager

import (
	"fmt"
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
			Name:         nodeName,
			TaskID:       task.ID,
			Dependencies: []string{},
		})
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

	var visit func(nodeName string, stack []string) error
	visit = func(nodeName string, stack []string) error {
		if state[nodeName] == 1 {
			return fmt.Errorf("%w: %s", ErrWorkflowCircularDependency, strings.Join(append(stack, nodeName), " -> "))
		}
		if state[nodeName] == 2 {
			return nil
		}
		state[nodeName] = 1

		for _, dep := range getDependencies(nodes, nodeName) {
			if err := visit(dep, append(stack, dep)); err != nil {
				return err
			}
		}
		state[nodeName] = 2
		return nil
	}

	for _, node := range nodes {
		if err := visit(node.TaskID, []string{}); err != nil {
			return err
		}
	}
	return nil
}

func getDependencies(nodes []model.Node, nodeName string) []string {
	for _, node := range nodes {
		if node.Name == nodeName {
			return node.Dependencies
		}
	}
	return []string{}
}
