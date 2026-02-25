package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kavos113/seseragi/internal/adapter/yaml"
)

func (c *Commands) AddTask(yamlPath string) error {
	absPath, err := filepath.Abs(yamlPath)
	if err != nil {
		return err
	}

	yamlData, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	task, err := yaml.LoadTaskInfoFromYAML(yamlData, absPath)
	if err != nil {
		return err
	}

	return c.tu.AddTask(*task)
}

func (c *Commands) ListTasks() error {
	tasks, err := c.tu.ListTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		fmt.Printf("%s - %s - %s\n", task.ID, task.Name, task.YamlPath)
	}

	return nil
}
