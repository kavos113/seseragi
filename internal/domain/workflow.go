package domain

import (
	"fmt"
	"time"
)

type Node struct {
	Name         string
	TaskName     string
	Dependencies []string // Nameのリスト
	Environments map[string]string
}

func GetNodeInputPath(nodeName string) string {
	return fmt.Sprintf("%s_input.json", nodeName)
}

func GetNodeOutputPath(nodeName string) string {
	return fmt.Sprintf("%s_output.json", nodeName)
}

type Workflow struct {
	ID          string
	Name        string
	RunInterval time.Duration
	Nodes       []Node
	YamlPath    string
}

type WorkflowRun struct {
	ID         string
	WorkflowID string
	StartTime  time.Time
	EndTime    time.Time
	Status     WorkflowStatus
}

type WorkflowStatus string

const (
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusCancelled WorkflowStatus = "cancelled"
)

type WorkflowRepository interface {
	CreateWorkflow(workflow Workflow) (Workflow, error)
	GetWorkflowByID(id string) (Workflow, error)
	GetAllWorkflows() ([]Workflow, error)

	UpdateWorkflow(workflow Workflow) (Workflow, error)
	AddNodeToWorkflow(workflowID string, node Node) (Workflow, error)
	DeleteNodeFromWorkflow(workflowID string, nodeName string) (Workflow, error)
	DeleteWorkflow(id string) error
}

type WorkflowRunRepository interface {
	CreateWorkflowRun(workflowRun WorkflowRun) (WorkflowRun, error)
	GetAllWorkflowRuns() ([]WorkflowRun, error)
	GetWorkflowRunByID(id string) (WorkflowRun, error)
	GetWorkflowRunsByWorkflowID(workflowID string) ([]WorkflowRun, error)
	GetWorkflowRunsAfter(workflowID string, before time.Time) ([]WorkflowRun, error)
}
