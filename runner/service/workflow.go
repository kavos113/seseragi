package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kavos113/seseragi/model"
	"github.com/kavos113/seseragi/model/repository/json"
)

type WorkflowRunner struct {
	workflowRepo    model.WorkflowRepository
	workflowRunRepo model.WorkflowRunRepository
	taskRepo        model.TaskRepository
}

func NewWorkflowRunner() *WorkflowRunner {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	jsonRepo := json.NewJsonRepository(filepath.Join(userConfigDir, "seseragi"))
	return &WorkflowRunner{
		workflowRepo:    json.NewJSONWorkflowRepository(jsonRepo),
		workflowRunRepo: json.NewJSONWorkflowRunRepository(jsonRepo),
		taskRepo:        json.NewJSONTaskRepository(jsonRepo),
	}
}

// interval: 1 hour
func (wr *WorkflowRunner) GetWorkflowToRun() ([]model.Workflow, error) {
	now := time.Now()
	limit := now.Add(-1 * time.Hour)

	workflows, err := wr.workflowRepo.GetAllWorkflows()
	if err != nil {
		return nil, err
	}

	var toRun []model.Workflow
	for _, workflow := range workflows {
		runs, err := wr.workflowRunRepo.GetWorkflowRunsAfter(workflow.ID, limit)
		if err != nil {
			return nil, err
		}

		if len(runs) == 0 {
			toRun = append(toRun, workflow)
		}
	}

	return toRun, nil
}

func (wr *WorkflowRunner) GetWorkflowByID(id string) (model.Workflow, error) {
	return wr.workflowRepo.GetWorkflowByID(id)
}

func (wr *WorkflowRunner) SaveWorkflowRun(run model.WorkflowRun) error {
	_, err := wr.workflowRunRepo.CreateWorkflowRun(run)
	return err
}

func (wr *WorkflowRunner) GetImageNameByTaskID(taskID string) (string, error) {
	task, err := wr.taskRepo.GetTaskByID(taskID)
	if err != nil {
		return "", err
	}
	return task.ImageName, nil
}

type NodeInfo struct {
	node      model.Node
	dependsOn []*NodeInfo
	doneCh    chan struct{} // 自分のタスクが完了したことを通知
	err       error
}

func (wr *WorkflowRunner) RunWorkflow(workflow model.Workflow, runNode func(model.Node) error) error {
	nodes := make(map[string]*NodeInfo)
	for _, node := range workflow.Nodes {
		nodes[node.TaskID] = &NodeInfo{
			node:   node,
			doneCh: make(chan struct{}),
			err:    nil,
		}
	}

	for _, nodeInfo := range nodes {
		for _, dep := range nodeInfo.node.Dependencies {
			depNodeInfo, ok := nodes[dep.TaskID]
			if !ok {
				return fmt.Errorf("dependency task %s not found", dep.TaskID)
			}
			nodeInfo.dependsOn = append(nodeInfo.dependsOn, depNodeInfo)
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(nodes))

	for _, nodeInfo := range nodes {
		go func(n *NodeInfo) {
			defer close(n.doneCh)
			defer wg.Done()

			for _, dep := range n.dependsOn {
				<-dep.doneCh

				if dep.err != nil {
					n.err = fmt.Errorf("dependency task %s failed: %w", dep.node.TaskID, dep.err)
					return
				}
			}
			fmt.Printf("all dependencies of task %s are completed, running task...\n", n.node.TaskID)
			n.err = runNode(n.node)
		}(nodeInfo)
	}

	wg.Wait()
	return nil
}
