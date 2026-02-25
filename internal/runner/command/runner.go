package command

import (
	"time"

	"github.com/kavos113/seseragi/internal/domain"
)

type CommandTaskRunner struct{
	Timeout time.Duration
}

func NewCommandTaskRunner(timeout time.Duration) *CommandTaskRunner {
	return &CommandTaskRunner{
		Timeout: timeout,
	}
}

func (r *CommandTaskRunner) Run(node domain.Node) error {

}
