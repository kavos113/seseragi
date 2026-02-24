package manager

import "errors"

var (
	ErrWorkflowCircularDependency = errors.New("circular dependency detected in workflow")
)
