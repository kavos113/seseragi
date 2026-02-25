package commands

import (
	"errors"
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

func (c *Commands) DeleteTask(taskID string) error {
	return c.tu.DeleteTask(taskID)
}

func (c *Commands) HandleTaskCommand(args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: seseragi task <subcommand>")
		fmt.Println("  Subcommands: add, list, delete")
		return errors.New("missing subcommand for task")
	}

	switch args[0] {
	case "add":
		if len(args) < 2 {
			fmt.Println("Usage: seseragi task add <yaml_path>")
			return errors.New("missing yaml path for adding task")
		}
		return c.AddTask(args[1])

	case "list":
		return c.ListTasks()

	case "delete":
		if len(args) < 2 {
			fmt.Println("Usage: seseragi task delete <task_id>")
			return errors.New("missing task ID for deleting task")
		}
		return c.DeleteTask(args[1])

	default:
		fmt.Printf("Unknown subcommand for task: %s\n", args[0])
		fmt.Println("Usage: seseragi task <subcommand>")
		fmt.Println("  Subcommands: add, list, delete")
		return fmt.Errorf("unknown subcommand for task: %s", args[0])
	}
}
