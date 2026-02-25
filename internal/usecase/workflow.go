package usecase

import "github.com/kavos113/seseragi/internal/domain"

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
	_, err := uc.workflowRepo.CreateWorkflow(workflow)
	return err
}

func (uc *workflowUseCase) ListWorkflows() ([]domain.Workflow, error) {
	return uc.workflowRepo.GetAllWorkflows()
}

func (uc *workflowUseCase) DeleteWorkflow(workflowID string) error {
	return uc.workflowRepo.DeleteWorkflow(workflowID)
}
