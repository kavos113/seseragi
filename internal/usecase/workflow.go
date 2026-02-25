package usecase

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kavos113/seseragi/internal/domain"
)

type WorkflowUseCase interface {
	AddWorkflow(workflow domain.Workflow) error
	ListWorkflows() ([]domain.Workflow, error)
	DeleteWorkflow(workflowID string) error
}

type workflowUseCase struct {
	workflowRepo domain.WorkflowRepository
}

func NewWorkflowUseCase(workflowRepo domain.WorkflowRepository) WorkflowUseCase {
	return &workflowUseCase{
		workflowRepo: workflowRepo,
	}
}

func (uc *workflowUseCase) AddWorkflow(workflow domain.Workflow) error {
	if err := checkCircularDependency(workflow.Nodes); err != nil {
		return err
	}
	if err := checkMissingDependency(workflow.Nodes); err != nil {
		return err
	}

	id := uuid.New().String()
	workflow.ID = id
	if _, err := uc.workflowRepo.CreateWorkflow(workflow); err != nil {
		return err
	}
	return nil
}

func (uc *workflowUseCase) ListWorkflows() ([]domain.Workflow, error) {
	return uc.workflowRepo.GetAllWorkflows()
}

func (uc *workflowUseCase) DeleteWorkflow(workflowID string) error {
	return uc.workflowRepo.DeleteWorkflow(workflowID)
}

func checkCircularDependency(nodes []domain.Node) error {
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

func checkMissingDependency(nodes []domain.Node) error {
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

func getDependencies(nodes []domain.Node, nodeName string) []string {
	for _, node := range nodes {
		if node.Name == nodeName {
			return node.Dependencies
		}
	}
	return []string{}
}
