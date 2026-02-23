package manager

import (
	"archive/tar"
	"bytes"
	"context"
	"os"

	"github.com/moby/moby/client"
)

type DockerClient struct {
	client *client.Client
}

func NewDockerClient() *DockerClient {
	cli, err := client.New(client.FromEnv, client.WithAPIVersionFromEnv())
	if err != nil {
		panic(err)
	}

	return &DockerClient{client: cli}
}

func (d *DockerClient) BuildImage(contextDir string, imageName string) error {
	var tarBuf bytes.Buffer
	tw := tar.NewWriter(&tarBuf)
	defer tw.Close()

	if err := tw.AddFS(os.DirFS(contextDir)); err != nil {
		return err
	}

	buildOptions := client.ImageBuildOptions{
		Tags:        []string{imageName},
		Dockerfile:  "Dockerfile",
		Remove:      true,
		ForceRemove: true,
	}

	resp, err := d.client.ImageBuild(context.Background(), &tarBuf, buildOptions)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
