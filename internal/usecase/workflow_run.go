package usecase

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kavos113/seseragi/internal/domain"
)

type WorkflowRunUseCase interface {
	RunWorkflow(workflowID string, runnerSelector func(node domain.Node) domain.NodeRunner) error
	GetWorkflowsToRun() ([]domain.Workflow, error)
	ListWorkflowRuns() ([]domain.WorkflowRun, error)
}

type workflowRunUseCase struct {
	workflowRepo    domain.WorkflowRepository
	workflowRunRepo domain.WorkflowRunRepository
	idGenerator     IDGenerator
}

func NewWorkflowRunUseCase(workflowRepo domain.WorkflowRepository, workflowRunRepo domain.WorkflowRunRepository, runner domain.NodeRunner, idGenerator IDGenerator) WorkflowRunUseCase {
	return &workflowRunUseCase{
		workflowRepo:    workflowRepo,
		workflowRunRepo: workflowRunRepo,
		idGenerator:     idGenerator,
	}
}

type NodeRun struct {
	node      domain.Node
	dependsOn []*NodeRun
	doneCh    chan struct{} // 自分が終了したことを通知
	runner    domain.NodeRunner
	err       error
}

func (uc *workflowRunUseCase) RunWorkflow(workflowID string, runnerSelector func(node domain.Node) domain.NodeRunner) error {
	workflow, err := uc.workflowRepo.GetWorkflowByID(workflowID)
	if err != nil {
		return err
	}

	nodes := make(map[string]*NodeRun)
	for _, node := range workflow.Nodes {
		nodes[node.Name] = &NodeRun{
			node:   node,
			doneCh: make(chan struct{}),
			runner: runnerSelector(node),
			err:    nil,
		}
	}

	for _, node := range workflow.Nodes {
		for _, depName := range node.Dependencies {
			depNode, exists := nodes[depName]
			if !exists {
				return errors.New("invalid workflow: dependency node not found: " + depName)
			}
			nodes[node.Name].dependsOn = append(nodes[node.Name].dependsOn, depNode)
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(nodes))

	for _, nodeRun := range nodes {
		go func(nr *NodeRun) {
			defer wg.Done()
			defer close(nr.doneCh)

			for _, dep := range nr.dependsOn {
				<-dep.doneCh

				if dep.err != nil {
					nr.err = errors.New("dependency failed: " + dep.node.Name)
					return
				}
			}

			fmt.Printf("Running node: %s\n", nr.node.Name)
			nr.err = nr.runner.Run(nr.node)
		}(nodeRun)
	}

	wg.Wait()
	for _, nodeRun := range nodes {
		if nodeRun.err != nil {
			return  fmt.Errorf("workflow run failed due to node error: %s error: %w", nodeRun.node.Name, nodeRun.err)
		}
	}

	return nil
}

// interval: 1 hour
func (uc *workflowRunUseCase) GetWorkflowsToRun() ([]domain.Workflow, error) {
	now := time.Now()
	limit := now.Add(-1 * time.Hour)

	workflows, err := uc.workflowRepo.GetAllWorkflows()
	if err != nil {
		return nil, err
	}

	var toRun []domain.Workflow
	for _, wf := range workflows {
		runs, err := uc.workflowRunRepo.GetWorkflowRunsAfter(wf.ID, limit)
		if err != nil {
			return nil, err
		}
		if len(runs) == 0 {
			toRun = append(toRun, wf)
		}
	}

	return toRun, nil
}

func (uc *workflowRunUseCase) ListWorkflowRuns() ([]domain.WorkflowRun, error) {
	return uc.workflowRunRepo.GetAllWorkflowRuns()
}
