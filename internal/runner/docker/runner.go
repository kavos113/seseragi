package docker

import (
	"errors"

	"github.com/kavos113/seseragi/internal/adapter/docker"
	"github.com/kavos113/seseragi/internal/domain"
)

type DockerNodeRunner struct {
	client *docker.Client
}

func NewDockerNodeRunner(client *docker.Client) *DockerNodeRunner {
	return &DockerNodeRunner{client: client}
}

func (r *DockerNodeRunner) Run(node domain.Node, task domain.Task, workflowRunID string) error {
	dockerDef, ok := task.TaskDef.(domain.DockerTaskDefinition)
	if !ok {
		return errors.New("invalid task definition type for DockerNodeRunner")
	}

	if err := r.client.RunContainer(dockerDef.ImageName, domain.GetDataDir(workflowRunID)); err != nil {
		return err
	}

	return nil
}
