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

	added, err := c.wu.AddWorkflow(*workflow)
	if err != nil {
		return err
	}
	fmt.Printf("Workflow '%s' added successfully with ID: %s\n", workflow.Name, added.ID)
	return nil
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

func (c *Commands) DeleteWorkflow(workflowID string) error {
	return c.wu.DeleteWorkflow(workflowID)
}

func (c *Commands) HandleWorkflowCommand(args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: seseragi workflow <subcommand>")
		fmt.Println("  Subcommands: add, list, delete")
		return fmt.Errorf("missing subcommand for workflow")
	}

	switch args[0] {
	case "add":
		if len(args) < 2 {
			fmt.Println("Usage: seseragi workflow add <yaml_path>")
			return fmt.Errorf("missing yaml path for adding workflow")
		}
		return c.AddWorkflow(args[1])

	case "list":
		return c.ListWorkflows()

	case "delete":
		if len(args) < 2 {
			fmt.Println("Usage: seseragi workflow delete <workflow_id>")
			return fmt.Errorf("missing workflow ID for deleting workflow")
		}
		return c.DeleteWorkflow(args[1])

	default:
		fmt.Printf("Unknown subcommand: %s\n", args[0])
		fmt.Println("Usage: seseragi workflow <subcommand>")
		fmt.Println("  Subcommands: add, list, delete")
		return fmt.Errorf("unknown subcommand for workflow: %s", args[0])
	}
}
