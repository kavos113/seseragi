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
		task, err := taskRepo.GetTaskByName(nodeInfo.Name)
		if err != nil {
			return model.Workflow{}, fmt.Errorf("failed to get task by name %s for node %s: %w", nodeInfo.Name, nodeName, err)
		}
		dependencies := make([]string, 0, len(nodeInfo.Dependencies))
		for _, dep := range nodeInfo.Dependencies {
			dependencies = append(dependencies, dep)
		}
		nodes = append(nodes, model.Node{
			Name:         nodeName,
			TaskID:       task.ID,
			Dependencies: dependencies,
		})
	}

	if err := checkMissingDependency(nodes); err != nil {
		return model.Workflow{}, err
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
		if err := visit(node.Name, []string{}); err != nil {
			return err
		}
	}
	return nil
}

func checkMissingDependency(nodes []model.Node) error {
	nodeNames := make(map[string]bool)
	for _, node := range nodes {
		nodeNames[node.Name] = true
	}

	for _, node := range nodes {
		for _, dep := range node.Dependencies {
			if !nodeNames[dep] {
				return fmt.Errorf("%w: node %s depends on missing node %s", ErrWorkflowMissingDependency, node.Name, dep)
			}
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
