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
	dp *docker.DockerTaskProvider
}

func NewCommands(
	tu usecase.TaskUseCase,
	wu usecase.WorkflowUseCase,
	wru usecase.WorkflowRunUseCase,
	dr *docker.DockerNodeRunner,
	dp *docker.DockerTaskProvider,
) *Commands {
	return &Commands{
		tu:  tu,
		wu:  wu,
		wru: wru,
		dr:  dr,
		dp:  dp,
	}
}
