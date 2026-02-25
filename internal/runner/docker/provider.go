package docker

import (
	"errors"

	"github.com/kavos113/seseragi/internal/adapter/docker"
	"github.com/kavos113/seseragi/internal/domain"
)

type DockerTaskProvider struct {
	client *docker.Client
}

func NewDockerTaskProvider(client *docker.Client) *DockerTaskProvider {
	return &DockerTaskProvider{client: client}
}

func (p *DockerTaskProvider) BuildTask(task domain.Task) error {
	dockerDef, ok := task.TaskDef.(domain.DockerTaskDefinition)
	if !ok {
		return errors.New("invalid task definition type for DockerTaskProvider")
	}

	if err := p.client.BuildImage(task.ID, dockerDef.ContextDir); err != nil {
		return err
	}

	return nil
}
