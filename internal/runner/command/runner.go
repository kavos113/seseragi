package command

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/kavos113/seseragi/internal/domain"
)

type CommandTaskRunner struct {
	Timeout time.Duration
}

func NewCommandTaskRunner(timeout time.Duration) *CommandTaskRunner {
	return &CommandTaskRunner{
		Timeout: timeout,
	}
}

func (r *CommandTaskRunner) Run(node domain.Node, task domain.Task) error {
	commandDef, ok := task.TaskDef.(domain.CommandTaskDefinition)
	if !ok {
		return nil // or return an error indicating invalid task definition
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "powershell", "-Command", commandDef.Command)
	case "linux", "darwin":
		cmd = exec.CommandContext(ctx, "sh", "-c", commandDef.Command)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	cmd.Dir = commandDef.WorkingDir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	fmt.Println(stdout.String())
	fmt.Println(stderr.String())

	if err != nil {
		return err
	}

	return nil
}
