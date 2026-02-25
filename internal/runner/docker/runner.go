package docker

import (
	"errors"

	"github.com/kavos113/seseragi/internal/adapter/docker"
	"github.com/kavos113/seseragi/internal/domain"
)

type DockerNodeRunner struct {
	client   *docker.Client
	taskRepo domain.TaskRepository
}

func NewDockerNodeRunner(client *docker.Client, taskRepo domain.TaskRepository) *DockerNodeRunner {
	return &DockerNodeRunner{client: client, taskRepo: taskRepo}
}

func (r *DockerNodeRunner) Run(node domain.Node) error {
	task, err := r.taskRepo.GetTaskByName(node.TaskName)
	if err != nil {
		return err
	}

	dockerDef, ok := task.TaskDef.(domain.DockerTaskDefinition)
	if !ok {
		return errors.New("invalid task definition type for DockerNodeRunner")
	}

	if err := r.client.RunContainer(dockerDef.ImageName); err != nil {
		return err
	}

	return nil
}
