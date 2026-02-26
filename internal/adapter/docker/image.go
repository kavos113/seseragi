package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/moby/moby/client"
)

func (c *Client) BuildImage(imageName string, contextDir string) error {
	var excludes []string

	ignorePath := filepath.Join(contextDir, ".dockerignore")
	f, err := os.Open(ignorePath)
	if err == nil {
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		line := strings.Split(string(content), "\n")
		for _, l := range line {
			l = strings.TrimSpace(l)
			if l != "" {
				excludes = append(excludes, l)
			}
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	var tarBuf bytes.Buffer
	tw := tar.NewWriter(&tarBuf)

	err = filepath.Walk(contextDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(contextDir, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		for _, pattern := range excludes {
			if strings.Contains(relPath, pattern) {
				return nil
			}
		}

		header := &tar.Header{
			Name:    relPath,
			Size:    info.Size(),
			Mode:    int64(info.Mode()),
			ModTime: info.ModTime(),
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(tw, file); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = tw.Close()
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
