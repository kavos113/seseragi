package service

import (
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
	jsonRepo := json.NewJsonRepository("data")
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