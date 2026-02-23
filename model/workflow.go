package model

import "time"

type Node struct {
	TaskID       string
	Dependencies []Task
}

type Workflow struct {
	ID          string
	Name        string
	RunInterval time.Duration
	Nodes       []Node
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

	UpdateWorkflow(workflow Workflow) (Workflow, error)
	AddNodeToWorkflow(workflowID string, node Node) (Workflow, error)
	DeleteNodeFromWorkflow(workflowID string, taskID string) (Workflow, error)
	DeleteWorkflow(id string) error
}

type WorkflowRunRepository interface {
	CreateWorkflowRun(workflowRun WorkflowRun) (WorkflowRun, error)
	GetWorkflowRunByID(id string) (WorkflowRun, error)
}
