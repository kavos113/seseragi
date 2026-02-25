package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/moby/moby/client"
)

func (c *Client) BuildImage(imageName string, contextDir string) error {
	var tarBuf bytes.Buffer
	tw := tar.NewWriter(&tarBuf)

	if err := tw.AddFS(os.DirFS(contextDir)); err != nil {
		return err
	}

	err := tw.Close()
	if err != nil {
		return err
	}

	buildOptions := client.ImageBuildOptions{
		Tags:        []string{imageName},
		Dockerfile:  "Dockerfile",
		Remove:      true,
		ForceRemove: true,
	}

	resp, err := c.client.ImageBuild(context.Background(), &tarBuf, buildOptions)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	for {
		var msg struct {
			Stream string `json:"stream"`
			Error  string `json:"error"`
		}
		if err := decoder.Decode(&msg); err != nil {
			break
		}
		if msg.Error != "" {
			return os.ErrInvalid
		}
		print(msg.Stream)
	}

	return nil
}
