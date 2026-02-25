package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kavos113/seseragi/cmd/cli/commands"
	"github.com/kavos113/seseragi/internal/adapter/docker"
	"github.com/kavos113/seseragi/internal/repository/json"
	dockerr "github.com/kavos113/seseragi/internal/runner/docker"
	"github.com/kavos113/seseragi/internal/usecase"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: seseragi <command> [args]")
		fmt.Println("  Available commands: task, workflow, run")
		os.Exit(1)
	}

	dc, err := docker.NewClient()
	if err != nil {
		panic(err)
	}

	appDataDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	jsonRepo := json.NewJsonRepository(filepath.Join(appDataDir, "seseragi"))
	taskRepo := json.NewJSONTaskRepository(jsonRepo)
	workflowRepo := json.NewJSONWorkflowRepository(jsonRepo)
	workflowRunRepo := json.NewJSONWorkflowRunRepository(jsonRepo)

	idGenerator := usecase.NewUUIDGenerator()

	dockerProvider := dockerr.NewDockerTaskProvider(dc)
	dockerRunner := dockerr.NewDockerNodeRunner(dc)

	taskUseCase := usecase.NewTaskUseCase(taskRepo, idGenerator)
	workflowUseCase := usecase.NewWorkflowUseCase(workflowRepo, taskRepo, idGenerator)
	workflowRunUseCase := usecase.NewWorkflowRunUseCase(workflowRepo, workflowRunRepo, taskRepo, idGenerator)

	cmds := commands.NewCommands(taskUseCase, workflowUseCase, workflowRunUseCase, dockerRunner, dockerProvider)

	switch os.Args[1] {
	case "task":
		err := cmds.HandleTaskCommand(os.Args[2:])
		if err != nil {
			fmt.Printf("Error handling task command: %v\n", err)
			os.Exit(1)
		}

	case "workflow":
		err := cmds.HandleWorkflowCommand(os.Args[2:])
		if err != nil {
			fmt.Printf("Error handling workflow command: %v\n", err)
			os.Exit(1)
		}

	case "run":
		err := cmds.HandleRunCommand(os.Args[2:])
		if err != nil {
			fmt.Printf("Error handling run command: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("  Available commands: task, workflow, run")
		os.Exit(1)
	}
}
