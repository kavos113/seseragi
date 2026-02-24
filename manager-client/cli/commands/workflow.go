package commands

import (
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
		println("ID:", wf.ID)
		println("Name:", wf.Name)
		println("YamlPath:", wf.YamlPath)
		println("Tasks:")
	}

	return nil
}