package service

import (
	"os"
	"path/filepath"
	"time"

	"github.com/kavos113/seseragi/model"
	"github.com/kavos113/seseragi/model/repository/json"
)

type WorkflowManager struct {
	workflowRepo    model.WorkflowRepository
	workflowRunRepo model.WorkflowRunRepository
}

func NewWorkflowRunner() *WorkflowManager {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	jsonRepo := json.NewJsonRepository(filepath.Join(userConfigDir, "seseragi"))
	return &WorkflowManager{
		workflowRepo:    json.NewJSONWorkflowRepository(jsonRepo),
		workflowRunRepo: json.NewJSONWorkflowRunRepository(jsonRepo),
	}
}

// interval: 1 hour
func (wm *WorkflowManager) GetWorkflowToRun() ([]model.Workflow, error) {
	now := time.Now()
	limit := now.Add(-1 * time.Hour)

	workflows, err := wm.workflowRepo.GetAllWorkflows()
	if err != nil {
		return nil, err
	}

	var toRun []model.Workflow
	for _, workflow := range workflows {
		runs, err := wm.workflowRunRepo.GetWorkflowRunsAfter(workflow.ID, limit)
		if err != nil {
			return nil, err
		}

		if len(runs) == 0 {
			toRun = append(toRun, workflow)
		}
	}

	return toRun, nil
}

func (wm *WorkflowManager) GetWorkflowByID(id string) (model.Workflow, error) {
	return wm.workflowRepo.GetWorkflowByID(id)
}

func (wm *WorkflowManager) SaveWorkflowRun(run model.WorkflowRun) error {
	_, err := wm.workflowRunRepo.CreateWorkflowRun(run)
	return err
}
