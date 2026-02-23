package commands

import (
	"path/filepath"

	"github.com/kavos113/seseragi/manager"
)

func BuildCommand(yamlPath string) error {
	absPath, err := filepath.Abs(yamlPath)
	if err != nil {
		return err
	}
	return manager.BuildTask(absPath)
}
