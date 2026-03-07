package docker

import (
	"bytes"
	"io"
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/kavos113/seseragi/internal/adapter/docker/mock_docker"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
)

func TestClient_RunContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMoby := mock_docker.NewMockmobyClient(ctrl)
	c := &Client{client: mockMoby}

	image := "test-image"
	dataDir := t.TempDir()
	nodeName := "test-node"
	envVars := map[string]string{"KEY": "VALUE"}

	mockMoby.EXPECT().ContainerCreate(gomock.Any(), gomock.Any()).Return(client.ContainerCreateResult{ID: "container-id"}, nil)
	
	waitRespCh := make(chan container.WaitResponse, 1)
	waitErrCh := make(chan error, 1)
	waitRespCh <- container.WaitResponse{StatusCode: 0}
	mockMoby.EXPECT().ContainerWait(gomock.Any(), "container-id", gomock.Any()).Return(client.ContainerWaitResult{
		Result: waitRespCh,
		Error:  waitErrCh,
	})
	
	mockMoby.EXPECT().ContainerStart(gomock.Any(), "container-id", gomock.Any()).Return(client.ContainerStartResult{}, nil)
	
	mockMoby.EXPECT().ContainerLogs(gomock.Any(), "container-id", gomock.Any()).Return(io.NopCloser(bytes.NewReader([]byte(""))), nil)
	
	mockMoby.EXPECT().ContainerRemove(gomock.Any(), "container-id", gomock.Any()).Return(client.ContainerRemoveResult{}, nil)

	err := c.RunContainer(image, dataDir, nodeName, envVars)
	assert.NoError(t, err)
}

func TestClient_BuildImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMoby := mock_docker.NewMockmobyClient(ctrl)
	c := &Client{client: mockMoby}

	imageName := "test-image"
	contextDir := t.TempDir()
	
	// Create a dummy Dockerfile
	// (BuildImage walks the contextDir, so we need some files)
	
	mockMoby.EXPECT().ImageBuild(gomock.Any(), gomock.Any(), gomock.Any()).Return(client.ImageBuildResult{
		Body: io.NopCloser(bytes.NewReader([]byte(`{"stream":"building..."}`))),
	}, nil)

	err := c.BuildImage(imageName, contextDir)
	assert.NoError(t, err)
}
