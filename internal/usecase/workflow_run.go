package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	taskRepo        domain.TaskRepository
	idGenerator     IDGenerator
}

func NewWorkflowRunUseCase(workflowRepo domain.WorkflowRepository, workflowRunRepo domain.WorkflowRunRepository, taskRepo domain.TaskRepository, idGenerator IDGenerator) WorkflowRunUseCase {
	return &workflowRunUseCase{
		workflowRepo:    workflowRepo,
		workflowRunRepo: workflowRunRepo,
		taskRepo:        taskRepo,
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

	run := domain.WorkflowRun{
		ID:         uc.idGenerator.GenerateID(),
		WorkflowID: workflow.ID,
		StartTime:  time.Now(),
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

	dataDir := domain.GetDataDir(run.ID)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	// defer os.RemoveAll(dataDir)

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

			if err := uc.collectOutputs(nr, run.ID); err != nil {
				nr.err = fmt.Errorf("failed to collect outputs for node %s: %w", nr.node.Name, err)
				return
			}

			fmt.Printf("Running node: %s\n", nr.node.Name)
			task, err := uc.taskRepo.GetTaskByName(nr.node.TaskName)
			if err != nil {
				nr.err = fmt.Errorf("failed to get task for node %s: %w", nr.node.Name, err)
				return
			}
			nr.err = nr.runner.Run(nr.node, task, run.ID)
		}(nodeRun)
	}

	wg.Wait()
	run.EndTime = time.Now()
	for _, nodeRun := range nodes {
		if nodeRun.err != nil {
			run.Status = domain.WorkflowStatusFailed
			if _, err := uc.workflowRunRepo.CreateWorkflowRun(run); err != nil {
				return err
			}
			return fmt.Errorf("workflow run failed due to node error: %s error: %w", nodeRun.node.Name, nodeRun.err)
		}
	}

	run.Status = domain.WorkflowStatusCompleted
	if _, err := uc.workflowRunRepo.CreateWorkflowRun(run); err != nil {
		return err
	}

	return nil
}

type NodeData struct {
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

type InputData []NodeData

// _output.json を結合して _input.json を作る. すべてのdependsOnが成功していることが前提
func (uc *workflowRunUseCase) collectOutputs(noderun *NodeRun, runID string) error {
	inputData := InputData{}

	for _, dep := range noderun.dependsOn {
		outputPath := filepath.Join(domain.GetDataDir(runID), domain.GetNodeOutputPath(dep.node.Name))

		outBytes, err := os.ReadFile(outputPath)
		if err != nil {
			return fmt.Errorf("failed to read output for node %s: %w", dep.node.Name, err)
		}

		inputData = append(inputData, NodeData{
			Name: dep.node.Name,
			Data: outBytes,
		})
	}

	inputBytes, err := json.Marshal(inputData)
	if err != nil {
		return fmt.Errorf("failed to marshal input data for node %s: %w", noderun.node.Name, err)
	}

	inputPath := filepath.Join(domain.GetDataDir(runID), domain.GetNodeInputPath(noderun.node.Name))
	if err := os.WriteFile(inputPath, inputBytes, 0644); err != nil {
		return fmt.Errorf("failed to write input file for node %s: %w", noderun.node.Name, err)
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
