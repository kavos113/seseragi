package docker

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/kavos113/seseragi/internal/adapter/docker/mock_docker"
	"github.com/kavos113/seseragi/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDockerNodeRunner_Run(t *testing.T) {
	tests := []struct {
		name          string
		node          domain.Node
		task          domain.Task
		workflowRunID string
		setupMock     func(m *mock_docker.MockDockerClient)
		wantErr       bool
	}{
		{
			name: "success",
			node: domain.Node{
				Name: "test-node",
				Environments: map[string]string{
					"ENV1": "VALUE1",
				},
			},
			task: domain.Task{
				ID:   "task-1",
				Name: "task-1",
				TaskDef: domain.DockerTaskDefinition{
					ImageName: "test-image",
				},
			},
			workflowRunID: "run-1",
			setupMock: func(m *mock_docker.MockDockerClient) {
				m.EXPECT().RunContainer("test-image", domain.GetDataDir("run-1"), "test-node", map[string]string{"ENV1": "VALUE1"}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failure: docker error",
			node: domain.Node{Name: "test-node"},
			task: domain.Task{
				Name: "task-1",
				TaskDef: domain.DockerTaskDefinition{
					ImageName: "test-image",
				},
			},
			workflowRunID: "run-1",
			setupMock: func(m *mock_docker.MockDockerClient) {
				m.EXPECT().RunContainer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("docker error"))
			},
			wantErr: true,
		},
		{
			name: "failure: wrong task definition",
			node: domain.Node{Name: "test-node"},
			task: domain.Task{
				Name:    "task-1",
				TaskDef: domain.CommandTaskDefinition{},
			},
			workflowRunID: "run-1",
			setupMock:     func(m *mock_docker.MockDockerClient) {},
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mock_docker.NewMockDockerClient(ctrl)
			tt.setupMock(mockClient)

			r := NewDockerNodeRunner(mockClient)
			err := r.Run(tt.node, tt.task, tt.workflowRunID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDockerTaskProvider_BuildTask(t *testing.T) {
	tests := []struct {
		name      string
		task      domain.Task
		setupMock func(m *mock_docker.MockDockerClient)
		wantErr   bool
	}{
		{
			name: "success",
			task: domain.Task{
				ID: "task-1",
				TaskDef: domain.DockerTaskDefinition{
					ContextDir: "/path/to/context",
				},
			},
			setupMock: func(m *mock_docker.MockDockerClient) {
				m.EXPECT().BuildImage("task-1", "/path/to/context").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failure: build error",
			task: domain.Task{
				ID: "task-1",
				TaskDef: domain.DockerTaskDefinition{
					ContextDir: "/path/to/context",
				},
			},
			setupMock: func(m *mock_docker.MockDockerClient) {
				m.EXPECT().BuildImage(gomock.Any(), gomock.Any()).Return(errors.New("build error"))
			},
			wantErr: true,
		},
		{
			name: "failure: wrong task definition",
			task: domain.Task{
				ID:      "task-1",
				TaskDef: domain.CommandTaskDefinition{},
			},
			setupMock: func(m *mock_docker.MockDockerClient) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mock_docker.NewMockDockerClient(ctrl)
			tt.setupMock(mockClient)

			p := NewDockerTaskProvider(mockClient)
			err := p.BuildTask(tt.task)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
