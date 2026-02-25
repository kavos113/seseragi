package domain

import (
	"fmt"
	"os"
	"path/filepath"
)

type NodeRunner interface {
	Run(node Node, task Task) error
}

func GetDataDir(workflowRunID string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("seseragi_data_%s", workflowRunID))
}