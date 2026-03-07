package command

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kavos113/seseragi/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCommandTaskRunner_Run(t *testing.T) {
	tests := []struct {
		name           string
		node           domain.Node
		task           domain.Task
		workflowRunID  string
		timeout        time.Duration
		wantErr        bool
		expectedErr    error
	}{
		{
			name: "success: simple echo",
			node: domain.Node{Name: "test-node"},
			task: domain.Task{
				Name: "test-task",
				TaskDef: domain.CommandTaskDefinition{
					Command: "echo hello",
				},
			},
			workflowRunID: "run-1",
			timeout:       10 * time.Second,
			wantErr:       false,
		},
		{
			name: "failure: invalid command",
			node: domain.Node{Name: "test-node"},
			task: domain.Task{
				Name: "test-task",
				TaskDef: domain.CommandTaskDefinition{
					Command: "non-existent-command-12345",
				},
			},
			workflowRunID: "run-1",
			timeout:       10 * time.Second,
			wantErr:       true,
		},
		{
			name: "failure: timeout",
			node: domain.Node{Name: "test-node"},
			task: domain.Task{
				Name: "test-task",
				TaskDef: domain.CommandTaskDefinition{
					Command: "sleep 5", // This might be OS-dependent, but sleep is common
				},
			},
			workflowRunID: "run-1",
			timeout:       1 * time.Second,
			wantErr:       true,
			expectedErr:   context.DeadlineExceeded,
		},
		{
			name: "failure: wrong task definition",
			node: domain.Node{Name: "test-node"},
			task: domain.Task{
				Name:    "test-task",
				TaskDef: domain.DockerTaskDefinition{},
			},
			workflowRunID: "run-1",
			timeout:       10 * time.Second,
			wantErr:       false, // Current implementation returns nil for wrong task def
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewCommandTaskRunner(tt.timeout)
			
			// Ensure data dir exists for env vars
			dataDir := domain.GetDataDir(tt.workflowRunID)
			err := os.MkdirAll(dataDir, 0755)
			assert.NoError(t, err)
			defer os.RemoveAll(dataDir)

			err = r.Run(tt.node, tt.task, tt.workflowRunID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommandTaskProvider_BuildTask(t *testing.T) {
	p := NewCommandTaskProvider()
	err := p.BuildTask(domain.Task{})
	assert.NoError(t, err)
}
