package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kavos113/seseragi/internal/adapter/yaml"
)

func (c *Commands) AddWorkflow(yamlPath string) error {
	absPath, err := filepath.Abs(yamlPath)
	if err != nil {
		return err
	}

	yamlData, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	workflow, err := yaml.LoadWorkflowInfoFromYAML(yamlData, absPath)
	if err != nil {
		return err
	}

	_, err = c.wu.AddWorkflow(*workflow)
	return err
}

func (c *Commands) ListWorkflows() error {
	workflows, err := c.wu.ListWorkflows()
	if err != nil {
		return err
	}

	for _, wf := range workflows {
		fmt.Printf("%s - %s - %s\n", wf.ID, wf.Name, wf.YamlPath)
	}

	return nil
}
