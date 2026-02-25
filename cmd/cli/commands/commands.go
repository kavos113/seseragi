package commands

import "github.com/kavos113/seseragi/internal/usecase"

type Commands struct {
	tu  usecase.TaskUseCase
	wu  usecase.WorkflowUseCase
	wru usecase.WorkflowRunUseCase
}

func NewCommands(tu usecase.TaskUseCase, wu usecase.WorkflowUseCase, wru usecase.WorkflowRunUseCase) *Commands {
	return &Commands{
		tu:  tu,
		wu:  wu,
		wru: wru,
	}
}
