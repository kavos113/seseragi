package commands

import (
	"github.com/kavos113/seseragi/internal/runner/docker"
	"github.com/kavos113/seseragi/internal/usecase"
)

type Commands struct {
	tu  usecase.TaskUseCase
	wu  usecase.WorkflowUseCase
	wru usecase.WorkflowRunUseCase

	dr *docker.DockerNodeRunner
}

func NewCommands(tu usecase.TaskUseCase, wu usecase.WorkflowUseCase, wru usecase.WorkflowRunUseCase, dr *docker.DockerNodeRunner) *Commands {
	return &Commands{
		tu:  tu,
		wu:  wu,
		wru: wru,
		dr:  dr,
	}
}
