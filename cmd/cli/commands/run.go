package commands

import (
	"fmt"

	"github.com/kavos113/seseragi/internal/domain"
)

func (c *Commands) RunWorkflow(workflowID string) error {
	runnerSelector := func(node domain.Node) domain.NodeRunner {
		// TODO: commandによる分岐
		return c.dr
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
