package main

import (
	"context"
	"errors"
	"slices"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type DockerClient struct {
	client  *client.Client
	running []string
}

func NewDockerClient() (*DockerClient, error) {
	cli, err := client.New(client.FromEnv, client.WithAPIVersionFromEnv())
	if err != nil {
		return nil, err
	}
	return &DockerClient{client: cli}, nil
}

func (dc *DockerClient) RunContainer(image string) error {
	ctx := context.Background()

	config := &container.Config{
		Image: image,
	}

	createOptions := client.ContainerCreateOptions{
		Config: config,
		Name:   image + "-container",
	}

	resp, err := dc.client.ContainerCreate(ctx, createOptions)
	if err != nil {
		return err
	}

	waitOptions := client.ContainerWaitOptions{
		Condition: container.WaitConditionNotRunning,
	}
	waitResp := dc.client.ContainerWait(ctx, resp.ID, waitOptions)

	startOptions := client.ContainerStartOptions{}
	if _, err := dc.client.ContainerStart(ctx, resp.ID, startOptions); err != nil {
		return err
	}

	dc.running = append(dc.running, resp.ID)

	select {
	case err := <-waitResp.Error:
		if err != nil {
			return err
		}

	case status := <-waitResp.Result:
		if status.StatusCode != 0 {
			return errors.New("container exited with non-zero status")
		}
		dc.running = slices.DeleteFunc(dc.running, func(id string) bool {
			return id == resp.ID
		})
	}

	return nil
}

func (dc *DockerClient) Cleanup() {
	ctx := context.Background()

	for _, containerID := range dc.running {
		timeout := 10
		stopOptions := client.ContainerStopOptions{
			Timeout: &timeout,
		}
		if _, err := dc.client.ContainerStop(ctx, containerID, stopOptions); err != nil {
			continue
		}

		removeOptions := client.ContainerRemoveOptions{
			Force: true,
		}
		if _, err := dc.client.ContainerRemove(ctx, containerID, removeOptions); err != nil {
			continue
		}
	}
}
