package manager

import "errors"

var (
	ErrWorkflowCircularDependency = errors.New("circular dependency detected in workflow")
	ErrWorkflowMissingDependency    = errors.New("missing dependency in workflow")
)
