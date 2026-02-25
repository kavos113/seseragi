package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func (c *Client) RunContainer(image string) error {
	ctx := context.Background()

	config := &container.Config{
		Image: image,
	}

	createOptions := client.ContainerCreateOptions{
		Config: config,
		Name:   image + "-container",
	}

	resp, err := c.client.ContainerCreate(ctx, createOptions)
	if err != nil {
		return err
	}

	waitOptions := client.ContainerWaitOptions{
		Condition: container.WaitConditionNotRunning,
	}
	waitResp := c.client.ContainerWait(ctx, resp.ID, waitOptions)

	startOptions := client.ContainerStartOptions{}
	if _, err := c.client.ContainerStart(ctx, resp.ID, startOptions); err != nil {
		return err
	}

	logsOptions := client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}
	stream, err := c.client.ContainerLogs(ctx, resp.ID, logsOptions)
	if err != nil {
		return err
	}
	defer stream.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	go func() {
		_, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, stream)
		if err != nil {
			return
		}

		fmt.Printf("Container %s logs:\nSTDOUT:\n%s\nSTDERR:\n%s\n", resp.ID, stdoutBuf.String(), stderrBuf.String())
	}()

	select {
	case err := <-waitResp.Error:
		if err != nil {
			return err
		}

	case status := <-waitResp.Result:
		if status.StatusCode != 0 {
			return errors.New("container exited with non-zero status")
		}

		removeOptions := client.ContainerRemoveOptions{
			Force: true,
		}
		if _, err := c.client.ContainerRemove(ctx, resp.ID, removeOptions); err != nil {
			return err
		}
	}

	return nil
}
