package commands

import (
	"fmt"
	"path/filepath"

	"github.com/kavos113/seseragi/manager"
)

func AddWorkflow(yamlPath string) error {
	absPath, err := filepath.Abs(yamlPath)
	if err != nil {
		return err
	}

	return manager.AddWorkflow(absPath)
}

func ShowWorkflows() error {
	workflows, err := manager.ListWorkflows()
	if err != nil {
		return err
	}

	for _, wf := range workflows {
		fmt.Printf("%s - %s - %s\n", wf.ID, wf.Name, wf.YamlPath)
	}

	return nil
}