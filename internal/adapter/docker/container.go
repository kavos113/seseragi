package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kavos113/seseragi/internal/domain"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func (c *Client) RunContainer(image string, dataDir string, nodeName string, envVars map[string]string) error {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	logDir := filepath.Join(confDir, "seseragi", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}
	logFilePath := filepath.Join(logDir, fmt.Sprintf("%s-%s.log", nodeName, time.Now().Format("20060102-150405")))
	logFile, err := os.Create(logFilePath)
	if err != nil {
		return err
	}
	defer logFile.Close()

	const containerDataDir = "/data"

	ctx := context.Background()

	envs := []string{}
	for key, value := range envVars {
		envs = append(envs, fmt.Sprintf("%s=%s", key, value))
	}
	envs = append(envs, fmt.Sprintf("WORKFLOW_INPUT_PATH=%s", fmt.Sprintf("%s/%s", containerDataDir, domain.GetNodeInputPath(nodeName))))
	envs = append(envs, fmt.Sprintf("WORKFLOW_OUTPUT_PATH=%s", fmt.Sprintf("%s/%s", containerDataDir, domain.GetNodeOutputPath(nodeName))))

	config := &container.Config{
		Image: image,
		Env:   envs,
	}
	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:%s", dataDir, containerDataDir),
		},
	}

	createOptions := client.ContainerCreateOptions{
		Config:     config,
		HostConfig: hostConfig,
	}

	resp, err := c.client.ContainerCreate(ctx, createOptions)
	if err != nil {
		return err
	}

	waitOptions := client.ContainerWaitOptions{
		Condition: container.WaitConditionNextExit,
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

		if _, err := logFile.WriteString(fmt.Sprintf("STDOUT:\n%s\nSTDERR:\n%s\n", stdoutBuf.String(), stderrBuf.String())); err != nil {
			return
		}
	}()

	select {
	case err := <-waitResp.Error:
		fmt.Printf("Error waiting for container %s: %v\n", resp.ID, err)
		if err != nil {
			return err
		}

	case status := <-waitResp.Result:
		fmt.Printf("Container %s exited with status code %d\n", resp.ID, status.StatusCode)
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

	fmt.Printf("Container %s finished\n", resp.ID)

	return nil
}
