package command

import "github.com/kavos113/seseragi/internal/domain"

type CommandTaskProvider struct{}

func NewCommandTaskProvider() *CommandTaskProvider {
	return &CommandTaskProvider{}
}

func (p *CommandTaskProvider) BuildTask(task domain.Task) error {
	// no need
	return nil
}
