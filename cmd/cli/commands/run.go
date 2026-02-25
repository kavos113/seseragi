package commands

import "github.com/kavos113/seseragi/internal/domain"

func (c *Commands) RunWorkflow(workflowID string) error {
	runnerSelector := func(node domain.Node) domain.NodeRunner {
		// TODO: commandによる分岐
		return c.dr
	}

	return c.wru.RunWorkflow(workflowID, runnerSelector)
}
