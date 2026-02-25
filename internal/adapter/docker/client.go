package docker

import (
	"github.com/moby/moby/client"
)

type Client struct {
	client *client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.New(client.FromEnv, client.WithAPIVersionFromEnv())
	if err != nil {
		return nil, err
	}
	return &Client{client: cli}, nil
}
