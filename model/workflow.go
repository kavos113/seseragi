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
}

type NodeRepository interface {
	CreateNode(node Node) (Node, error)
	GetNodeByID(id string) (Node, error)
	
	UpdateNode(node Node) (Node, error)
	DeleteNode(id string) error
}

type WorkflowRepository interface {
	CreateWorkflow(workflow Workflow) (Workflow, error)
	GetWorkflowByID(id string) (Workflow, error)
	
	UpdateWorkflow(workflow Workflow) (Workflow, error)
	AddNodeToWorkflow(workflowID string, node Node) (Workflow, error)
	DeleteWorkflow(id string) error
}

type WorkflowRunRepository interface {
	CreateWorkflowRun(workflowRun WorkflowRun) (WorkflowRun, error)
	GetWorkflowRunByID(id string) (WorkflowRun, error)
	
	UpdateWorkflowRun(workflowRun WorkflowRun) (WorkflowRun, error)
	DeleteWorkflowRun(id string) error
}