package commands

import (
	"os"
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/kavos113/seseragi/internal/domain"
	"github.com/kavos113/seseragi/internal/usecase/mock_usecase"
	"github.com/stretchr/testify/assert"
)

func TestAddTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskUC := mock_usecase.NewMockTaskUseCase(ctrl)
	c := &Commands{tu: mockTaskUC}

	// Create dummy task YAML
	yamlContent := `
name: test-task
type: command
command: echo hello
`
	tmpFile, err := os.CreateTemp("", "task-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	
	_, err = tmpFile.WriteString(yamlContent)
	assert.NoError(t, err)
	tmpFile.Close()

	mockTaskUC.EXPECT().AddTask(gomock.Any(), gomock.Any()).Return(nil)

	err = c.AddTask(tmpFile.Name())
	assert.NoError(t, err)
}

func TestAddWorkflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWorkflowUC := mock_usecase.NewMockWorkflowUseCase(ctrl)
	c := &Commands{wu: mockWorkflowUC}

	// Create dummy workflow YAML
	yamlContent := `
name: test-workflow
trigger:
  type: interval
  interval: 60
nodes:
  node-1:
    name: task-1
`
	tmpFile, err := os.CreateTemp("", "workflow-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	
	_, err = tmpFile.WriteString(yamlContent)
	assert.NoError(t, err)
	tmpFile.Close()

	mockWorkflowUC.EXPECT().AddWorkflow(gomock.Any()).Return(domain.Workflow{ID: "wf-1"}, nil)

	err = c.AddWorkflow(tmpFile.Name())
	assert.NoError(t, err)
}

func TestRunWorkflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWorkflowRunUC := mock_usecase.NewMockWorkflowRunUseCase(ctrl)
	c := &Commands{wru: mockWorkflowRunUC}

	mockWorkflowRunUC.EXPECT().RunWorkflow("wf-1", gomock.Any()).Return(nil)

	err := c.RunWorkflow("wf-1")
	assert.NoError(t, err)
}
