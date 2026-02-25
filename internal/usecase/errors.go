package usecase

import "errors"

var (
	ErrWorkflowCircularDependency = errors.New("circular dependency detected in workflow")
	ErrWorkflowMissingDependency  = errors.New("missing dependency in workflow")
	ErrWorkflowMissingTask        = errors.New("workflow contains node with missing task")
)
