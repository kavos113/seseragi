package commands

import (
	"fmt"
	"time"

	"github.com/kavos113/seseragi/internal/domain"
	"github.com/kavos113/seseragi/internal/runner/command"
)

func (c *Commands) RunWorkflow(workflowID string) error {
	runnerSelector := func(node domain.Node) domain.NodeRunner {
		t, err := c.wu.GetTaskTypeFromNode(node)
		if err != nil {
			fmt.Printf("Error determining task type for node %s: %v\n", node.Name, err)
			return nil
		}

		switch t {
		case domain.TaskTypeDocker:
			return c.dr

		case domain.TaskTypeCommand:
			return command.NewCommandTaskRunner(5 * time.Minute)

		default:
			fmt.Printf("No runner available for task type %s in node %s\n", t, node.Name)
			return nil
		}
	}

	return c.wru.RunWorkflow(workflowID, runnerSelector)
}

func (c *Commands) HandleRunCommand(args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: seseragi run <workflow_id>")
		return fmt.Errorf("missing workflow ID to run")
	}

	workflowID := args[0]
	return c.RunWorkflow(workflowID)
}
